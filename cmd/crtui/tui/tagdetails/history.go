package tagdetails

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/samber/lo"
)

func (m *Model) drawHistory() string {
	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			MarginTop(1).
			Render("❯ HISTORY"),
	}

	for i, step := range m.tag.Platforms[m.activeTabIndex].Config.History {
		stepElems := []string{}

		if step.CreatedBy != "" {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(12).
						Render("Created By:"),
					lipgloss.NewStyle().
						Width(lo.Ternary(step.EmptyLayer, m.platformVP.Width-29, m.platformVP.Width-16)).
						Render(step.CreatedBy),
					lo.Ternary(step.EmptyLayer,
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("#000000")).
							Background(lipgloss.Color("#FFFFFF")).
							Bold(true).
							Render(" EMPTY LAYER "),
						"",
					),
				),
			)
		}

		if !step.Created.IsZero() {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(12).
						Render("Created:"),
					lipgloss.NewStyle().
						Width(m.platformVP.Width-16).
						Render(step.Created.String()),
				),
			)
		}

		if step.Comment != "" {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(12).
						Render("Comment:"),
					lipgloss.NewStyle().
						Width(m.platformVP.Width-16).
						Render(step.Comment),
				),
			)
		}

		if step.Author != "" {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(12).
						Render("Author:"),
					lipgloss.NewStyle().
						Width(m.platformVP.Width-16).
						Render(step.Author),
				),
			)
		}

		strs = append(strs,
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.NewStyle().
					Faint(true).
					Bold(true).
					Render(fmt.Sprintf("#%d", i+1)),
				" ",
				lipgloss.JoinVertical(lipgloss.Top,
					stepElems...,
				),
			),
			" ",
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
