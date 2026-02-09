package ui

import (
	"cmp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type WindowConfig struct {
	Width, Height                                int
	LeftTitle, RightTitle                        string
	Content                                      string
	LeftTitleColor, RightTitleColor, BorderColor lipgloss.Color
}

func NewWindow(cfg WindowConfig) string {
	leftTitleColor := cmp.Or(cfg.LeftTitleColor, lipgloss.Color("#959595"))
	rightTitleColor := cmp.Or(cfg.RightTitleColor, lipgloss.Color("#959595"))
	borderColor := cmp.Or(cfg.BorderColor, lipgloss.Color("#383838"))

	border := lipgloss.RoundedBorder()

	// LEFT TITLE
	var renderedLeftTitle string
	leftTitleWidth := 0

	if cfg.LeftTitle != "" {
		leftTitleStyle := lipgloss.NewStyle().
			Foreground(leftTitleColor).
			Bold(true).
			Padding(0, 1)
		renderedLeftTitle = leftTitleStyle.Render(cfg.LeftTitle)
		leftTitleWidth = lipgloss.Width(renderedLeftTitle)
	}

	// RIGHT TITLE
	var renderedRightTitle string
	rightTitleWidth := 0

	if cfg.RightTitle != "" {
		rightTitleStyle := lipgloss.NewStyle().
			Foreground(rightTitleColor).
			Padding(0, 1)
		renderedRightTitle = rightTitleStyle.Render(cfg.RightTitle)
		rightTitleWidth = lipgloss.Width(renderedRightTitle)
	}

	// BORDERS
	prefix := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(border.TopLeft + border.Top)

	suffixWidth := max(cfg.Width-leftTitleWidth-rightTitleWidth-1, 0)

	middleLine := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat(border.Top, suffixWidth))

	topRightChar := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(border.TopRight)

	topBar := prefix + renderedLeftTitle + middleLine + renderedRightTitle + topRightChar

	// BODY
	bodyStyle := lipgloss.NewStyle().
		Border(border, false, true, true, true).
		BorderForeground(borderColor).
		Width(cfg.Width).
		Height(cfg.Height - 1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		bodyStyle.Render(cfg.Content),
	)
}
