package tagdetails

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

func (m *Model) drawCmd() string {
	cmds := m.tag.Platforms[m.activeTabIndex].Config.Cmd

	if len(cmds) == 0 {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Render("❯ CMD"),
	}

	for _, cmd := range cmds {
		strs = append(strs,
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.NewStyle().
					Faint(true).
					Render("● "),
				lipgloss.NewStyle().
					Render(cmd),
			),
		)
	}

	strs = append(strs, "")

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
