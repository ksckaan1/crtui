package tagdetails

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *Model) drawEnvs(activePlatform models.Platform, width int) string {
	envs := activePlatform.Config.Env

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
			m.drawField(split[0], split[1], width),
			"",
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
