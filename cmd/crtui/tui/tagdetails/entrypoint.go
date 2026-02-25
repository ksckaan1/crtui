package tagdetails

import (
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *Model) drawEntrypoint(activePlatform models.Platform, width int) string {
	entrypoint := activePlatform.Config.Entrypoint

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
					Width(width-2).
					Render(ep),
			),
		)
	}

	strs = append(strs, "")

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
