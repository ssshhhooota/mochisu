package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	var content string
	switch m.screen {
	case screenRepoSelect:
		dialog := repoSelectStyle.Render(m.form.View())
		content = lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dialog,
		)
	case screenIssueList:
		paneWidth := m.width / 2
		issueListStyle := paneBaseStyle.
			Width(paneWidth - 2).
			Height(m.height - 2)
		previewStyle := paneBaseStyle.
			Width(paneWidth - 2).
			Height(m.height - 2)
		if m.focus == focusList {
			issueListStyle = issueListStyle.BorderForeground(colorAccent)
			previewStyle = previewStyle.BorderForeground(colorDim)
		} else {
			issueListStyle = issueListStyle.BorderForeground(colorDim)
			previewStyle = previewStyle.BorderForeground(colorAccent)
		}
		issueList := issueListStyle.Render(m.issues.View())
		preview := previewStyle.Render(m.preview.View())
		content = lipgloss.JoinHorizontal(lipgloss.Top, issueList, preview)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func renderMarkdown(src string, width int) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(styles.TokyoNightStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return src
	}
	out, err := r.Render(src)
	if err != nil {
		return src
	}
	return out
}
