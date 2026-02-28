package tagdetails

import (
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *TagDetailsScreenModel) drawWorkingDir(activePlatform models.Platform, width int) string {
	wd := activePlatform.Config.WorkingDir

	if wd == "" {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Width(width).
			Render("❯ WORKING DIR"),
		lipgloss.NewStyle().
			MarginBottom(1).
			Width(width).
			Render(wd),
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
