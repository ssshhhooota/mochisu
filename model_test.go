package main

import "testing"

func TestIssueListItemInterface(t *testing.T) {
	tests := []struct {
		name       string
		iss        issue
		wantTitle  string
		wantDesc   string
		wantFilter string
	}{
		{
			name:       "regular issue",
			iss:        issue{Name: "Fix bug", Body: "details", Number: 1, URL: "u"},
			wantTitle:  "Fix bug",
			wantDesc:   "",
			wantFilter: "Fix bug",
		},
		{
			name:       "zero value",
			iss:        issue{},
			wantTitle:  "",
			wantDesc:   "",
			wantFilter: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.iss.Title(); got != tc.wantTitle {
				t.Errorf("Title() = %q, want %q", got, tc.wantTitle)
			}
			if got := tc.iss.Description(); got != tc.wantDesc {
				t.Errorf("Description() = %q, want %q", got, tc.wantDesc)
			}
			if got := tc.iss.FilterValue(); got != tc.wantFilter {
				t.Errorf("FilterValue() = %q, want %q", got, tc.wantFilter)
			}
		})
	}
}

func TestNewModel(t *testing.T) {
	repos := []repo{
		{NameWithOwner: "owner/repo1"},
		{NameWithOwner: "owner/repo2"},
	}
	m := newModel(repos)

	if m.form == nil {
		t.Fatal("form is nil")
	}
	if len(m.repos) != len(repos) {
		t.Errorf("repos length = %d, want %d", len(m.repos), len(repos))
	}
	if m.screen != screenRepoSelect {
		t.Errorf("screen = %v, want screenRepoSelect", m.screen)
	}
	if m.focus != focusList {
		t.Errorf("focus = %v, want focusList", m.focus)
	}
	if m.selectedRepo != "" {
		t.Errorf("selectedRepo = %q, want empty", m.selectedRepo)
	}
	if m.selectedIssue != (issue{}) {
		t.Errorf("selectedIssue = %+v, want zero value", m.selectedIssue)
	}
}

func TestNewRepoForm(t *testing.T) {
	tests := []struct {
		name  string
		repos []repo
	}{
		{"nil repos", nil},
		{"empty repos", []repo{}},
		{"single repo", []repo{{NameWithOwner: "a/b"}}},
		{"multiple repos", []repo{{NameWithOwner: "a/b"}, {NameWithOwner: "c/d"}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("newRepoForm panicked: %v", r)
				}
			}()
			f := newRepoForm(tc.repos)
			if f == nil {
				t.Fatal("form is nil")
			}
		})
	}
}
