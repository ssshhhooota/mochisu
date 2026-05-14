package main

import (
	"log/slog"

	"charm.land/bubbles/v2/list"
	"charm.land/huh/v2"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global: keys that work on any screen
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.screen == screenIssueList {
				if m.focus == focusList {
					m.focus = focusPreview
				} else {
					m.focus = focusList
				}
				return m, nil
			}
		case "H":
			if m.screen == screenIssueList &&
				m.issues.FilterState() != list.Filtering {
				m.focus = focusList
				return m, nil
			}
		case "L":
			if m.screen == screenIssueList &&
				m.issues.FilterState() != list.Filtering {
				m.focus = focusPreview
				return m, nil
			}
		case "o":
			if m.screen == screenIssueList &&
				m.focus == focusList &&
				m.issues.FilterState() != list.Filtering &&
				m.selectedIssue.Number > 0 {
				repo := m.selectedRepo
				number := m.selectedIssue.Number
				return m, func() tea.Msg {
					if err := openIssueInBrowser(repo, number); err != nil {
						slog.Error("failed to open issue in browser", "err", err)
					}
					return nil
				}
			}
		case "ctrl+r":
			if m.screen == screenIssueList {
				m.screen = screenRepoSelect
				m.selectedRepo = ""
				m.selectedIssue = issue{}
				m.focus = focusList
				m.issues.SetItems(nil)
				m.preview.SetContent("")
				m.form = newRepoForm(m.repos)
				if m.width > 0 {
					m.form = m.form.WithWidth(m.width * 3 / 4)
				}
				return m, m.form.Init()
			}
		}
	}

	// Global: window size
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = ws.Width, ws.Height
		paneWidth := ws.Width / 2
		m.issues.SetSize(paneWidth-2, ws.Height-2)
		m.preview.SetWidth(paneWidth - 2)
		m.preview.SetHeight(ws.Height - 2)
		m.form = m.form.WithWidth(ws.Width * 3 / 4)
		m.refreshPreview()
	}

	switch m.screen {
	case screenRepoSelect:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			m.selectedRepo = m.form.GetString("repo")
			m.screen = screenIssueList
			issues, err := fetchIssues(m.selectedRepo)
			if err != nil {
				slog.Error("failed to fetch issues", "err", err)
				return m, cmd
			}
			items := make([]list.Item, 0, len(issues))
			for _, iss := range issues {
				items = append(items, iss)
			}
			m.issues.SetItems(items)
			if sel, ok := m.issues.SelectedItem().(issue); ok {
				m.selectedIssue = sel
			}
			m.refreshPreview()
		}
		return m, cmd

	case screenIssueList:
		var cmd tea.Cmd
		prevTitle := m.selectedIssue.Name
		if m.focus == focusList {
			m.issues, cmd = m.issues.Update(msg)
			if sel, ok := m.issues.SelectedItem().(issue); ok {
				m.selectedIssue = sel
			}
			if m.selectedIssue.Name != prevTitle {
				m.refreshPreview()
				m.preview.GotoTop()
			}
		} else {
			m.preview, cmd = m.preview.Update(msg)
		}
		return m, cmd
	}

	return m, nil
}

func (m *model) refreshPreview() {
	width := m.preview.Width()
	if width <= 0 {
		return
	}
	m.preview.SetContent(renderMarkdown(m.selectedIssue.Body, width))
}
