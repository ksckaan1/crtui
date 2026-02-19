package tagdetails

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

func (m *Model) drawUser() string {
	user := m.tag.Platforms[m.activeTabIndex].Config.User

	if user == "" {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Width(m.platformVP.Width).
			Render("❯ USER"),
		lipgloss.NewStyle().
			MarginBottom(1).
			Width(m.platformVP.Width).
			Render(user),
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
