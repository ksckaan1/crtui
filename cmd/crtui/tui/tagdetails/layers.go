package tagdetails

import (
	"fmt"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/samber/lo"
)

func (m *TagDetailsScreenModel) drawLayers(activePlatform models.Platform, width int) string {
	sizeBytes := lo.SumBy(
		activePlatform.Layers,
		func(item models.Layer) int {
			return item.Size
		},
	)

	strs := []string{
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().
				Foreground(ui.PrimaryColor).
				Bold(true).
				Width(width/2).
				Render("❯ LAYERS"),
			lipgloss.NewStyle().
				Bold(true).
				AlignHorizontal(lipgloss.Right).
				Width(width/2).
				Render(fmt.Sprintf("%s (%dB)", humanReadableSize(sizeBytes), sizeBytes)),
		),
	}

	layers := []string{
		lipgloss.NewStyle().
			Foreground(ui.SecondaryColor).
			Render(fmt.Sprintf("╭%s╮", strings.Repeat("─", width-2))),
	}

	layerCount := len(activePlatform.Layers)

	for i, layer := range activePlatform.Layers {
		size := humanReadableSize(layer.Size)

		if size != fmt.Sprintf("%dB", layer.Size) {
			size += fmt.Sprintf(" (%dB)", layer.Size)
		}

		out := []string{
			m.drawPercent(sizeBytes, layer.Size, width-2),
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().
					Faint(true).
					Width(9).
					Render("Digest:"),
				lipgloss.NewStyle().
					Width(width-11).
					Render(layer.Digest),
			),
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().
					Faint(true).
					Width(9).
					Render("Size:"),
				lipgloss.NewStyle().
					Width(width-11).
					Render(size),
			),
		}

		layers = append(layers,
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(ui.SecondaryColor).
				BorderTop(false).
				BorderBottom(false).
				Render(lipgloss.JoinVertical(lipgloss.Top, out...)),
		)

		if i != layerCount-1 {
			layers = append(layers, lipgloss.NewStyle().
				Foreground(ui.SecondaryColor).
				Render(fmt.Sprintf("├%s┤", strings.Repeat("─", width-2))))
		}
	}

	if layerCount > 0 {
		layers = append(layers,
			lipgloss.NewStyle().
				Foreground(ui.SecondaryColor).
				Render(fmt.Sprintf("╰%s╯", strings.Repeat("─", width-2))),
		)
	}

	strs = append(strs, layers...)

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}

func (m *TagDetailsScreenModel) drawPercent(total, current int, width int) string {
	if total <= 0 {
		return strings.Repeat("░", width-8) + "   0.0%"
	}

	rawPercent := (float64(current) / float64(total)) * 100

	displayPercent := math.Floor(rawPercent*10) / 10

	if current < total && displayPercent >= 100.0 {
		displayPercent = 99.9
	}

	if current >= total {
		displayPercent = 100.0
	}

	barWidth := max(width-8, 0)

	filledLength := int(math.Floor(float64(barWidth) * (displayPercent / 100)))
	emptyLength := barWidth - filledLength

	bar := strings.Repeat("█", filledLength) + strings.Repeat("░", emptyLength)

	return fmt.Sprintf("%s %6.1f%%", bar, displayPercent)
}
