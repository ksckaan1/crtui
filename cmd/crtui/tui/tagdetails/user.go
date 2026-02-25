package tagdetails

import (
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *Model) drawUser(activePlatform models.Platform, width int) string {
	user := activePlatform.Config.User

	if user == "" {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Width(width).
			Render("❯ USER"),
		lipgloss.NewStyle().
			MarginBottom(1).
			Width(width).
			Render(user),
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
