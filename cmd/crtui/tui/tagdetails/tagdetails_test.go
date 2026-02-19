package tagdetails

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

func TestTab(t *testing.T) {
	// ─ │ ╭ ╮ ╰ ╯ ┼ ├ ┤ ┬ ┴

	tabs := []string{
		"Example Tab 1",
		"Example Tab 2",
		"Example Tab 3",
		"Example Tab 4",
	}

	activeTab := 1

	width := 100
	left := 0
	overflow := false

	drawnTabs := []string{}
	for i, title := range tabs {
		drawnTabs = append(drawnTabs, drawTab(title, i == activeTab))
	}

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		drawnTabs...,
	)

	overflow = lipgloss.Width(out) > width

	fmt.Println(lipgloss.JoinHorizontal(
		lipgloss.Top,
		"\n\n╭",
		lo.Ternary(left > 0, "«\n«\n«", ""),
		cropHorizontal(out, left, width),
		lo.Ternary(overflow, "»\n»\n»", ""),
		"\n\n╮",
	))
}

//╭───────────╮ ╭───────────╮
//│  Active   │ │ Inactive  │
//╯           ╰─┴───────────┴────────────

func drawTab(title string, active bool) string {
	pLine := strings.Repeat("─", len(title))
	pEmpty := strings.Repeat(" ", len(title))

	sb := &strings.Builder{}

	fmt.Fprint(sb, "╭─")
	fmt.Fprint(sb, pLine)
	fmt.Fprint(sb, "─╮\n")
	fmt.Fprint(sb, "│ ")
	fmt.Fprint(sb, title)
	fmt.Fprint(sb, " │\n")

	if active {
		fmt.Fprint(sb, "╯ ")
		fmt.Fprint(sb, pEmpty)
		fmt.Fprint(sb, " ╰")
	} else {
		fmt.Fprint(sb, "┴")
		fmt.Fprint(sb, pLine)
		fmt.Fprint(sb, "──┴")
	}

	return sb.String()
}

func cropHorizontal(s string, left, length int) string {
	if length == 0 {
		return ""
	}

	left = max(left, 0)

	lines := strings.Split(s, "\n")

	for i, line := range lines {
		runesLine := []rune(line)

		right := min(left+length, len(runesLine))

		if right <= left {
			lines[i] = ""
			continue
		}

		fmt.Println(left, right, len(runesLine))

		lines[i] = string(runesLine[left:right])
	}

	return strings.Join(lines, "\n")
}
