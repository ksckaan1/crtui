package tagdetails

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

func (m *Model) drawRootFS() string {
	rootFS := m.tag.Platforms[m.activeTabIndex].Config.RootFS

	strs := []string{
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().
				Foreground(ui.PrimaryColor).
				Bold(true).
				Width(m.platformVP.Width/2).
				Render("❯ ROOTFS"),
			lipgloss.NewStyle().
				Bold(true).
				AlignHorizontal(lipgloss.Right).
				Width(m.platformVP.Width/2).
				Render(lipgloss.JoinHorizontal(lipgloss.Top,
					"Type: ",
					rootFS.Type,
				)),
		),
		"",
	}

	for i, diffID := range rootFS.DiffIDs {
		strs = append(strs,
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.NewStyle().
					Faint(true).
					Render(fmt.Sprintf("#%03d ", i+1)),
				lipgloss.NewStyle().
					Width(m.platformVP.Width-3).
					Render(diffID),
			),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
