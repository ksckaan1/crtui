package ui

import (
	"cmp"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

type WindowConfig struct {
	Width, Height                                int
	LeftTitle, RightTitle                        string
	Content                                      string
	LeftTitleColor, RightTitleColor, BorderColor color.Color
}

func NewWindow(cfg WindowConfig) string {
	leftTitleColor := cmp.Or(cfg.LeftTitleColor, lipgloss.Color("#959595"))
	rightTitleColor := cmp.Or(cfg.RightTitleColor, lipgloss.Color("#959595"))
	borderColor := cmp.Or(cfg.BorderColor, SecondaryColor)

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

	topRightChar := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(border.TopRight)

	suffixWidth := max(cfg.Width-leftTitleWidth-rightTitleWidth-3, 0) // 3 = prefix + top right char

	middleLine := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat(border.Top, suffixWidth))

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
