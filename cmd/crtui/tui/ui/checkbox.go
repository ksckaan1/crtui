package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/samber/lo"
)

type Checkbox struct {
	checked   bool
	isFocused bool
}

func NewCheckbox() Checkbox {
	return Checkbox{}
}

func (c *Checkbox) Value() bool {
	return c.checked
}

func (c *Checkbox) SetValue(value bool) {
	c.checked = value
}

func (c *Checkbox) Toggle() {
	c.checked = !c.checked
}

func (c *Checkbox) Focused() bool {
	return c.isFocused
}

func (c *Checkbox) Focus() tea.Cmd {
	c.isFocused = true
	return nil
}

func (c *Checkbox) Blur() {
	c.isFocused = false
}

func (c *Checkbox) Reset() {
	c.checked = false
}

func (c *Checkbox) View() string {
	body := lo.Ternary(c.checked, "[X]", "[ ]")

	style := lipgloss.NewStyle()

	if c.isFocused {
		style = style.
			Background(lipgloss.White).
			Foreground(lipgloss.Black)
	}

	return style.Render(body)
}

func (c *Checkbox) Update(msg tea.Msg) (Checkbox, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case (msg.String() == "space" || msg.String() == "enter") && c.isFocused:
			c.Toggle()
		}
	}

	return *c, nil
}
