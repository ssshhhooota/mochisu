package main

import "charm.land/lipgloss/v2"

var (
	colorAccent      = lipgloss.Color("63")
	colorDim         = lipgloss.Color("240")
	colorStateOpen   = lipgloss.Color("#3FB950")
	colorStateClosed = lipgloss.Color("#8957E5")
)

var (
	repoSelectStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 2)

	paneBaseStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder())
)
