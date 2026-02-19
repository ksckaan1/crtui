package tagdetails

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

func (m *Model) drawEnvs() string {
	envs := m.tag.Platforms[m.activeTabIndex].Config.Env

	if len(envs) == 0 {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Render("❯ ENVIRONMENT VARIABLES"),
	}

	for _, env := range envs {
		split := strings.SplitN(env, "=", 2)
		if len(split) != 2 {
			continue
		}

		strs = append(strs,
			m.drawField(split[0], split[1], m.platformVP.Width-2),
			"",
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
