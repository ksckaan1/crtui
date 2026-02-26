package ui

import (
	"cmp"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

type Window struct {
	width, height   int
	leftTitle       string
	leftTitleColor  color.Color
	rightTitle      string
	rightTitleColor color.Color
	borderColor     color.Color
	content         string
}

func NewWindow() *Window {
	return &Window{}
}

func (m *Window) SetWidth(width int) {
	m.width = width
}

func (m *Window) Width() int {
	return m.width
}

func (m *Window) SetHeight(height int) {
	m.height = height
}

func (m *Window) Height() int {
	return m.height
}

func (m *Window) SetLeftTitle(s string) {
	m.leftTitle = s
}

func (m *Window) LeftTitle() string {
	return m.leftTitle
}

func (m *Window) SetRightTitle(s string) {
	m.rightTitle = s
}

func (m *Window) RightTitle() string {
	return m.rightTitle
}

func (m *Window) SetLeftTitleColor(c color.Color) {
	m.leftTitleColor = c
}

func (m *Window) SetRightTitleColor(c color.Color) {
	m.rightTitleColor = c
}

func (m *Window) SetBorderColor(c color.Color) {
	m.borderColor = c
}

func (m *Window) SetContent(cb func(width, height int) string) {
	m.content = cb(m.width-1, m.height-1)
}

func (m *Window) View() string {
	leftTitleColor := cmp.Or(m.leftTitleColor, lipgloss.Color("#959595"))
	rightTitleColor := cmp.Or(m.rightTitleColor, lipgloss.Color("#959595"))
	borderColor := cmp.Or(m.borderColor, SecondaryColor)

	border := lipgloss.RoundedBorder()

	// LEFT TITLE
	var renderedLeftTitle string
	leftTitleWidth := 0

	if m.leftTitle != "" {
		leftTitleStyle := lipgloss.NewStyle().
			Foreground(leftTitleColor).
			Bold(true).
			Padding(0, 1)
		renderedLeftTitle = leftTitleStyle.Render(m.leftTitle)
		leftTitleWidth = lipgloss.Width(renderedLeftTitle)
	}

	// RIGHT TITLE
	var renderedRightTitle string
	rightTitleWidth := 0

	if m.rightTitle != "" {
		rightTitleStyle := lipgloss.NewStyle().
			Foreground(rightTitleColor).
			Padding(0, 1)
		renderedRightTitle = rightTitleStyle.Render(m.rightTitle)
		rightTitleWidth = lipgloss.Width(renderedRightTitle)
	}

	// BORDERS
	prefix := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(border.TopLeft + border.Top)

	topRightChar := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(border.TopRight)

	suffixWidth := max(m.width-leftTitleWidth-rightTitleWidth-3, 0) // 3 = prefix + top right char

	middleLine := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat(border.Top, suffixWidth))

	topBar := prefix + renderedLeftTitle + middleLine + renderedRightTitle + topRightChar

	// BODY
	bodyStyle := lipgloss.NewStyle().
		Border(border, false, true, true, true).
		BorderForeground(borderColor).
		Width(m.width).
		Height(m.height - 1)

	vp := viewport.New(
		viewport.WithWidth(m.width-1),
		viewport.WithHeight(m.height-1),
	)

	vp.SetContent(m.content)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		bodyStyle.Render(vp.View()),
	)
}
