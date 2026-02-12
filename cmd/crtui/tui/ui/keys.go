package ui

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type KeysWindowConfig struct {
	Width, Height int
	Keys          []key.Binding
}

func NewKeysWindow(cfg KeysWindowConfig) string {
	w := NewWindow(WindowConfig{
		Width:     cfg.Width,
		Height:    cfg.Height,
		LeftTitle: "Keys",
		Content:   drawKeysContent(cfg.Keys, cfg.Height),
	})

	return w
}

var (
	keysContainer = lipgloss.NewStyle().
			PaddingLeft(1)
	keysKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(6)
	keysDescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262"))
)

func drawKeysContent(keys []key.Binding, height int) string {
	cols := make([]string, 0, int(math.Ceil(float64(len(keys))/float64(height))))

	var col string
	for i, k := range keys {
		col += fmt.Sprintf("%s %s", keysKeyStyle.Render(k.Help().Key), keysDescriptionStyle.Render(k.Help().Desc))
		if i%height != 3 {
			col += "\n"
			continue
		}

		cols = append(cols, col, "  ")
		col = ""
	}

	if col != "" {
		cols = append(cols, col)
	}

	keysContent := lipgloss.JoinHorizontal(lipgloss.Top, cols...)

	return keysContainer.Render(keysContent)
}
