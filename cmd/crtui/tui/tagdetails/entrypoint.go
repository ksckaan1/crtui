package tagdetails

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

func (m *Model) drawEntrypoint() string {
	entrypoint := m.tag.Platforms[m.activeTabIndex].Config.Entrypoint

	if len(entrypoint) == 0 {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Render("❯ ENTRYPOINT"),
	}

	for _, ep := range entrypoint {
		strs = append(strs,
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.NewStyle().
					Faint(true).
					Render("● "),
				lipgloss.NewStyle().
					Render(ep),
			),
		)
	}

	strs = append(strs, "")

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
