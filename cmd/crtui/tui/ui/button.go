package ui

import (
	"image/color"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Button struct {
	width     int
	isFocused bool
	text      string

	normalForegroundColor  color.Color
	normalBackgroundColor  color.Color
	focusedForegroundColor color.Color
	focusedBackgroundColor color.Color
}

func NewButton(text string) Button {
	return Button{
		text: text,

		normalForegroundColor:  lipgloss.Black,
		normalBackgroundColor:  PrimaryColor,
		focusedForegroundColor: lipgloss.Black,
		focusedBackgroundColor: lipgloss.White,
	}
}

func (c *Button) SetNormalColors(fg color.Color, bg color.Color) {
	c.normalForegroundColor = fg
	c.normalBackgroundColor = bg
}

func (c *Button) SetFocusedColors(fg color.Color, bg color.Color) {
	c.focusedForegroundColor = fg
	c.focusedBackgroundColor = bg
}

func (c *Button) Width() int {
	return c.width
}

func (c *Button) SetWidth(width int) {
	c.width = width
}

func (c *Button) Focused() bool {
	return c.isFocused
}

func (c *Button) Focus() tea.Cmd {
	c.isFocused = true
	return nil
}

func (c *Button) Blur() {
	c.isFocused = false
}

func (c *Button) Reset() {

}

func (c *Button) View() string {
	style := lipgloss.NewStyle().
		Width(c.width).
		Bold(true).
		Align(lipgloss.Center).
		Background(c.normalBackgroundColor).
		Foreground(c.normalForegroundColor)

	if c.isFocused {
		style = style.
			Width(c.width).
			Bold(true).
			Align(lipgloss.Center).
			Background(c.focusedBackgroundColor).
			Foreground(c.focusedForegroundColor)
	}

	return style.Render(c.text)
}
