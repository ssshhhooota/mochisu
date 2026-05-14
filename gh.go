package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cli/go-gh/v2"
)

var ghExec = gh.ExecContext

func ghJSON[T any](timeout time.Duration, args ...string) (T, error) {
	var zero T
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stdout, stderr, err := ghExec(ctx, args...)
	if err != nil {
		return zero, fmt.Errorf("gh %v failed: %w: %s", args, err, stderr.String())
	}

	var result T
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return zero, fmt.Errorf("failed to parse: %w", err)
	}
	return result, nil
}

func fetchRepos(owner string) ([]repo, error) {
	args := []string{"repo", "list", "--json", "nameWithOwner"}
	if owner != "" {
		args = append(args, owner)
	}
	return ghJSON[[]repo](3*time.Second, args...)
}

// gh issue list --json -h
func fetchIssues(repo string) ([]issue, error) {
	return ghJSON[[]issue](3*time.Second, "issue", "list", "--json", "title,body,number,url", "-R", repo)
}

func openIssueInBrowser(repo string, number int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, stderr, err := ghExec(ctx, "issue", "view", strconv.Itoa(number), "-R", repo, "--web")
	if err != nil {
		return fmt.Errorf("gh issue view --web failed: %w: %s", err, stderr.String())
	}
	return nil
}
