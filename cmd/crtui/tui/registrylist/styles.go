package registrylist

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.Color("#7D56F4")

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
