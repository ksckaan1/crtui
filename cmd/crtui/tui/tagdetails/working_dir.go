package tagdetails

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

func (m *Model) drawWorkingDir() string {
	wd := m.tag.Platforms[m.activeTabIndex].Config.WorkingDir

	if wd == "" {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Width(m.platformVP.Width).
			Render("❯ WORKING DIR"),
		lipgloss.NewStyle().
			MarginBottom(1).
			Width(m.platformVP.Width).
			Render(wd),
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
