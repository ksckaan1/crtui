package figlet

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/version"
)

var Figlet = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Faint(false).Render(FigletText)

const figletText = `                       ╭─╮
╭────╮╭────╮╭─╮  ╭─╮╭─╮├─┤
│ ╭──╯│ ╭─┬┴╯ ╰─╮│ ││ ││ │
│ │   │ │ ╰─╮ ╭─╯│ ││ ││ │
│ ╰──╮│ │   │ ╰─╮│ ╰╯ ││ │
╰────╯╰─╯   ╰───╯╰────╯╰─╯`

var FigletText = func() string {
	versionText := lipgloss.NewStyle().
		Foreground(lipgloss.Black).
		Background(lipgloss.Color("#444444")).
		Padding(0, 1).
		Bold(true).
		Render(version.Version)

	backLines := strings.Split(figletText, "\n")

	backHeight := len(backLines)
	backWidth := 0
	for _, l := range backLines {
		if w := lipgloss.Width(l); w > backWidth {
			backWidth = w
		}
	}

	frontLines := strings.Split(versionText, "\n")
	frontHeight := len(frontLines)

	frontWidth := 0
	for _, l := range frontLines {
		if w := lipgloss.Width(l); w > frontWidth {
			frontWidth = w
		}
	}

	startY := max(backHeight-frontHeight, 0)

	startX := max(backWidth-frontWidth, 0)

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
				leftPart = string(runes[:end])
			}

			rightPart := ""
			rightStart := startX + lipgloss.Width(fLine)
			if rightStart < totalRunes {
				rightPart = string(runes[rightStart:])
			}

			finalView.WriteString(leftPart + fLine + rightPart + "\n")
		} else {
			finalView.WriteString(line + "\n")
		}
	}

	return strings.TrimSuffix(finalView.String(), "\n")
}()
