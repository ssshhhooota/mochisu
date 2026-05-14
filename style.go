package main

import "charm.land/lipgloss/v2"

var (
	colorAccent = lipgloss.Color("63")
	colorDim    = lipgloss.Color("240")
)

var (
	repoSelectStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 2)

	paneBaseStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder())
)
