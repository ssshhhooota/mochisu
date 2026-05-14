package main

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func updateModel(t *testing.T, m model, msg tea.Msg) (model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	mm, ok := next.(model)
	if !ok {
		t.Fatalf("Update returned %T, want model", next)
	}
	return mm, cmd
}

func TestUpdateQuitKeys(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
	}{
		{"q", tea.KeyPressMsg{Code: 'q'}},
		{"ctrl+c", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newModel(nil)
			_, cmd := updateModel(t, m, tc.msg)
			if cmd == nil {
				t.Fatal("expected non-nil cmd")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Errorf("cmd produced %T, want tea.QuitMsg", cmd())
			}
		})
	}
}

func TestUpdateTabTogglesFocus(t *testing.T) {
	m := newModel(nil)
	m.screen = screenIssueList
	m.focus = focusList

	tab := tea.KeyPressMsg{Code: tea.KeyTab}

	m, _ = updateModel(t, m, tab)
	if m.focus != focusPreview {
		t.Errorf("after first tab: focus = %v, want focusPreview", m.focus)
	}
	m, _ = updateModel(t, m, tab)
	if m.focus != focusList {
		t.Errorf("after second tab: focus = %v, want focusList", m.focus)
	}
}

func TestUpdateTabIgnoredOnRepoSelect(t *testing.T) {
	m := newModel(nil)
	// initial screen is screenRepoSelect, focus is focusList
	got, _ := updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyTab})
	if got.focus != focusList {
		t.Errorf("focus changed on repo-select: got %v, want focusList", got.focus)
	}
}

func TestUpdateHLKeys(t *testing.T) {
	tests := []struct {
		name       string
		key        tea.KeyPressMsg
		startFocus focus
		wantFocus  focus
	}{
		{"H switches to list", tea.KeyPressMsg{Code: 'H'}, focusPreview, focusList},
		{"L switches to preview", tea.KeyPressMsg{Code: 'L'}, focusList, focusPreview},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := newModel(nil)
			m.screen = screenIssueList
			m.focus = tc.startFocus
			got, _ := updateModel(t, m, tc.key)
			if got.focus != tc.wantFocus {
				t.Errorf("focus = %v, want %v", got.focus, tc.wantFocus)
			}
		})
	}
}

func TestUpdateOpenInBrowserInvokesGh(t *testing.T) {
	var capturedArgs []string
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		capturedArgs = args
		return bytes.Buffer{}, bytes.Buffer{}, nil
	})()

	m := newModel(nil)
	m.screen = screenIssueList
	m.focus = focusList
	m.selectedRepo = "owner/repo"
	m.selectedIssue = issue{Number: 42, Name: "test"}

	_, cmd := updateModel(t, m, tea.KeyPressMsg{Code: 'o'})
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}
	cmd()
	wantArgs := []string{"issue", "view", "42", "-R", "owner/repo", "--web"}
	if !reflect.DeepEqual(capturedArgs, wantArgs) {
		t.Errorf("captured args = %v, want %v", capturedArgs, wantArgs)
	}
}

func TestUpdateOpenInBrowserSkippedWhenFocusPreview(t *testing.T) {
	called := false
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		called = true
		return bytes.Buffer{}, bytes.Buffer{}, nil
	})()

	m := newModel(nil)
	m.screen = screenIssueList
	m.focus = focusPreview
	m.selectedRepo = "owner/repo"
	m.selectedIssue = issue{Number: 42, Name: "test"}

	_, cmd := updateModel(t, m, tea.KeyPressMsg{Code: 'o'})
	if cmd != nil {
		cmd()
	}
	if called {
		t.Error("openIssueInBrowser should not be invoked when focus is preview")
	}
}

func TestUpdateOpenInBrowserSkippedWhenNoIssue(t *testing.T) {
	called := false
	defer withGhExec(func(ctx context.Context, args ...string) (bytes.Buffer, bytes.Buffer, error) {
		called = true
		return bytes.Buffer{}, bytes.Buffer{}, nil
	})()

	m := newModel(nil)
	m.screen = screenIssueList
	m.focus = focusList
	// selectedIssue.Number is zero — no issue selected.

	_, cmd := updateModel(t, m, tea.KeyPressMsg{Code: 'o'})
	if cmd != nil {
		cmd()
	}
	if called {
		t.Error("openIssueInBrowser should not be invoked when no issue selected")
	}
}

func TestUpdateCtrlRResetsState(t *testing.T) {
	repos := []repo{{NameWithOwner: "owner/repo"}}
	m := newModel(repos)
	m.screen = screenIssueList
	m.selectedRepo = "owner/repo"
	m.selectedIssue = issue{Number: 1, Name: "x"}
	m.focus = focusPreview
	m.preview.SetContent("previous content")

	ctrlR := tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl}
	got, _ := updateModel(t, m, ctrlR)

	if got.screen != screenRepoSelect {
		t.Errorf("screen = %v, want screenRepoSelect", got.screen)
	}
	if got.selectedRepo != "" {
		t.Errorf("selectedRepo = %q, want empty", got.selectedRepo)
	}
	if got.selectedIssue != (issue{}) {
		t.Errorf("selectedIssue = %+v, want zero value", got.selectedIssue)
	}
	if got.focus != focusList {
		t.Errorf("focus = %v, want focusList", got.focus)
	}
	if items := got.issues.Items(); len(items) != 0 {
		t.Errorf("issues items length = %d, want 0", len(items))
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := newModel(nil)
	got, _ := updateModel(t, m, tea.WindowSizeMsg{Width: 200, Height: 50})

	if got.width != 200 {
		t.Errorf("width = %d, want 200", got.width)
	}
	if got.height != 50 {
		t.Errorf("height = %d, want 50", got.height)
	}
	// paneWidth = 200 / 2 = 100; preview width = paneWidth - 2 = 98
	if w := got.preview.Width(); w != 98 {
		t.Errorf("preview width = %d, want 98", w)
	}
	// preview height = 50 - 2 = 48
	if h := got.preview.Height(); h != 48 {
		t.Errorf("preview height = %d, want 48", h)
	}
}
