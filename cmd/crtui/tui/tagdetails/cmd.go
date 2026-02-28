package tagdetails

import (
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *TagDetailsScreenModel) drawCmd(activePlatform models.Platform, width int) string {
	cmds := activePlatform.Config.Cmd

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
					Width(width-2).
					Render(cmd),
			),
		)
	}

	strs = append(strs, "")

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
