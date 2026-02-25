package tagdetails

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *Model) drawRootFS(activePlatform models.Platform, width int) string {
	rootFS := activePlatform.Config.RootFS

	strs := []string{
		"",
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().
				Foreground(ui.PrimaryColor).
				Bold(true).
				Width(width/2).
				Render("❯ ROOTFS"),
			lipgloss.NewStyle().
				Bold(true).
				AlignHorizontal(lipgloss.Right).
				Width(width/2).
				Render(lipgloss.JoinHorizontal(lipgloss.Top,
					"Type: ",
					rootFS.Type,
				)),
		),
		lipgloss.NewStyle().
			Foreground(ui.SecondaryColor).
			Render(fmt.Sprintf("╭%s╮", strings.Repeat("─", width-2))),
	}

	for i, diffID := range rootFS.DiffIDs {
		strs = append(strs,
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), false, true, false, true).
				BorderForeground(ui.SecondaryColor).
				Render(
					lipgloss.JoinHorizontal(lipgloss.Top,
						lipgloss.NewStyle().
							Faint(true).
							Bold(true).
							Render(fmt.Sprintf("#%03d ", i+1)),
						lipgloss.NewStyle().
							Width(width-7).
							Render(diffID),
					),
				),
		)

		if i != len(rootFS.DiffIDs)-1 {
			strs = append(strs,
				lipgloss.NewStyle().
					Foreground(ui.SecondaryColor).
					Render(fmt.Sprintf("├%s┤", strings.Repeat("─", width-2))),
			)
		}
	}

	strs = append(strs,
		lipgloss.NewStyle().
			Foreground(ui.SecondaryColor).
			Render(fmt.Sprintf("╰%s╯", strings.Repeat("─", width-2))),
	)

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
