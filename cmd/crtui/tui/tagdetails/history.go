package tagdetails

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
)

func (m *TagDetailsScreenModel) drawHistory(activePlatform models.Platform, width int) string {
	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginTop(1).
			Render("❯ HISTORY"),
		lipgloss.NewStyle().
			Foreground(ui.SecondaryColor).
			Render(fmt.Sprintf("╭%s╮", strings.Repeat("─", width-2))),
	}

	var (
		keyWidth   = 12
		valueWidth = width - 17
	)

	history := activePlatform.Config.History

	for i, step := range history {
		stepElems := []string{}

		var emptyLayerBadge string

		if step.EmptyLayer {
			emptyLayerBadge = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Render(" EMPTY LAYER ")
		}

		if step.CreatedBy != "" {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(keyWidth).
						Render("Created By:"),
					lipgloss.NewStyle().
						Width(valueWidth-lipgloss.Width(emptyLayerBadge)-2).
						Render(step.CreatedBy),
					emptyLayerBadge,
				),
			)
		}

		if !step.Created.IsZero() {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(keyWidth).
						Render("Created:"),
					lipgloss.NewStyle().
						Width(valueWidth).
						Render(step.Created.String()),
				),
			)
		}

		if step.Comment != "" {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(keyWidth).
						Render("Comment:"),
					lipgloss.NewStyle().
						Width(valueWidth).
						Render(step.Comment),
				),
			)
		}

		if step.Author != "" {
			stepElems = append(stepElems,
				lipgloss.JoinHorizontal(lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Width(keyWidth).
						Render("Author:"),
					lipgloss.NewStyle().
						Width(valueWidth).
						Render(step.Author),
				),
			)
		}

		strs = append(strs,
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), false, true, false, true).
				BorderForeground(ui.SecondaryColor).
				Width(width).
				Render(
					lipgloss.JoinHorizontal(lipgloss.Top,
						lipgloss.NewStyle().
							Faint(true).
							Bold(true).
							Render(fmt.Sprintf("#%03d ", i+1)),
						lipgloss.JoinVertical(lipgloss.Top,
							stepElems...,
						),
					),
				),
		)

		if i != len(history)-1 {
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
