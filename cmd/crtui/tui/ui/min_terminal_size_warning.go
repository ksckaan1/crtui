package ui

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

type MinTerminalSizeWarning struct {
	minWidth, minHeight int
}

func NewMinTerminalSizeWarning(width, height int) *MinTerminalSizeWarning {
	return &MinTerminalSizeWarning{
		minWidth:  width,
		minHeight: height,
	}
}

func (w *MinTerminalSizeWarning) View() string {
	currentWidth, currentHeight, err := term.GetSize(os.Stdin.Fd())
	if err != nil {
		return "error when getting terminal size"
	}

	if currentWidth < w.minWidth || currentHeight < w.minHeight {

		return lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Width(currentWidth).
			Height(currentHeight).
			Render(fmt.Sprintf("Minimum terminal size must be %dx%d", w.minWidth, w.minHeight))
	}

	return ""
}
