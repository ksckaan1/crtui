package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func NewOverlay(status *Status) *Overlay {
	return &Overlay{
		status: status,
	}
}

type Overlay struct {
	back   string
	front  string
	status *Status
}

func (o *Overlay) SetBack(back string) {
	o.back = back
}

func (o *Overlay) SetFront(front string) {
	o.front = front
}

func (o Overlay) View() string {
	cleanBack := string(ansi.Strip(o.back))
	backLines := strings.Split(cleanBack, "\n")

	faintStyle := lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#444444"))

	backHeight := len(backLines)
	backWidth := 0
	for _, l := range backLines {
		if w := lipgloss.Width(l); w > backWidth {
			backWidth = w
		}
	}

	frontLines := strings.Split(o.front, "\n")
	frontHeight := len(frontLines)

	frontWidth := 0
	for _, l := range frontLines {
		if w := lipgloss.Width(l); w > frontWidth {
			frontWidth = w
		}
	}

	startY := (backHeight - frontHeight) / 2
	startX := (backWidth - frontWidth) / 2

	var finalView strings.Builder
	for y, line := range backLines {
		if y >= startY && y < startY+frontHeight {
			fLineIdx := y - startY
			fLine := frontLines[fLineIdx]

			runes := []rune(line)
			totalRunes := len(runes)

			leftPart := ""
			if startX > 0 {
				end := min(startX, totalRunes)
				leftPart = faintStyle.Render(string(runes[:end]))
			}

			rightPart := ""
			rightStart := startX + lipgloss.Width(fLine)
			if rightStart < totalRunes {
				rightPart = faintStyle.Render(string(runes[rightStart:]))
			}

			finalView.WriteString(leftPart + fLine + rightPart + "\n")
		} else {
			finalView.WriteString(faintStyle.Render(line) + "\n")
		}
	}

	out := strings.TrimSuffix(finalView.String(), "\n")

	lines := strings.Split(out, "\n")

	out = strings.Join(lines[:len(lines)-1-o.status.Height()], "\n")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		out,
		fmt.Sprintf(" %s", o.status.View()),
	)
}
