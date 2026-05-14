package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
)

type screen int

const (
	screenRepoSelect screen = iota
	screenIssueList
)

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func initialModel(owner string) (model, error) {
	repos, err := fetchRepos(owner)
	if err != nil {
		return model{}, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	return newModel(repos), nil
}

func main() {
	// Logger
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				slog.Error("failed to close log file", "error", err)
			}
		}()
		slog.SetDefault(slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}
	var owner string
	if len(os.Args) > 1 {
		owner = os.Args[1]
	}
	m, err := initialModel(owner)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
