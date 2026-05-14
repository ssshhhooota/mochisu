package main

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

type ghExecFn = func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error)

// withGhExec swaps ghExec for the duration of a test and returns a restore fn.
func withGhExec(fn ghExecFn) func() {
	orig := ghExec
	ghExec = fn
	return func() { ghExec = orig }
}

func bufOf(s string) bytes.Buffer { return *bytes.NewBufferString(s) }

func TestGhJSONSuccess(t *testing.T) {
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		return bufOf(`[{"nameWithOwner":"a/b"}]`), bytes.Buffer{}, nil
	})()

	got, err := ghJSON[[]repo](time.Second, "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []repo{{NameWithOwner: "a/b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGhJSONCommandError(t *testing.T) {
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		return bytes.Buffer{}, bufOf("auth required"), errors.New("exit status 1")
	})()

	_, err := ghJSON[[]repo](time.Second, "x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGhJSONInvalidJSON(t *testing.T) {
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		return bufOf("not json"), bytes.Buffer{}, nil
	})()

	_, err := ghJSON[[]repo](time.Second, "x")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestFetchReposArgs(t *testing.T) {
	tests := []struct {
		name     string
		owner    string
		wantArgs []string
	}{
		{"no owner", "", []string{"repo", "list", "--json", "nameWithOwner"}},
		{"with owner", "octocat", []string{"repo", "list", "--json", "nameWithOwner", "octocat"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var capturedArgs []string
			defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
				capturedArgs = args
				return bufOf(`[]`), bytes.Buffer{}, nil
			})()

			if _, err := fetchRepos(tc.owner); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(capturedArgs, tc.wantArgs) {
				t.Errorf("args = %v, want %v", capturedArgs, tc.wantArgs)
			}
		})
	}
}

func TestFetchIssues(t *testing.T) {
	var capturedArgs []string
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		capturedArgs = args
		return bufOf(`[{"title":"t","body":"b","number":1,"url":"u"}]`), bytes.Buffer{}, nil
	})()

	got, err := fetchIssues("owner/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantArgs := []string{"issue", "list", "--json", "title,body,number,url", "-R", "owner/repo"}
	if !reflect.DeepEqual(capturedArgs, wantArgs) {
		t.Errorf("args = %v, want %v", capturedArgs, wantArgs)
	}
	want := []issue{{Name: "t", Body: "b", Number: 1, URL: "u"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestOpenIssueInBrowserArgs(t *testing.T) {
	var capturedArgs []string
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		capturedArgs = args
		return bytes.Buffer{}, bytes.Buffer{}, nil
	})()

	if err := openIssueInBrowser("owner/repo", 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantArgs := []string{"issue", "view", "42", "-R", "owner/repo", "--web"}
	if !reflect.DeepEqual(capturedArgs, wantArgs) {
		t.Errorf("args = %v, want %v", capturedArgs, wantArgs)
	}
}

func TestOpenIssueInBrowserError(t *testing.T) {
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		return bytes.Buffer{}, bufOf("network down"), errors.New("exit status 1")
	})()

	if err := openIssueInBrowser("owner/repo", 1); err == nil {
		t.Fatal("expected error, got nil")
	}
}
