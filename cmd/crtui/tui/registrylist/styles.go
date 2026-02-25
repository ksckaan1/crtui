package registrylist

import "charm.land/lipgloss/v2"

var (
	onlineDot = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render("●")

	offlineDot = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render("●")

	unauthDot = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Render("●")
)
