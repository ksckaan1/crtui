package tagdetails

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/samber/lo"
)

func (m *Model) drawLayers() string {
	sizeBytes := lo.SumBy(
		m.tag.Platforms[m.activeTabIndex].Layers,
		func(item models.Layer) int {
			return item.Size
		},
	)

	strs := []string{
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().
				Foreground(ui.PrimaryColor).
				Bold(true).
				Width(m.platformVP.Width/2).
				Render("❯ LAYERS"),
			lipgloss.NewStyle().
				Bold(true).
				AlignHorizontal(lipgloss.Right).
				Width(m.platformVP.Width/2).
				Render(fmt.Sprintf("%s (%dB)", humanReadableSize(sizeBytes), sizeBytes)),
		),
	}

	layers := []string{
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#383838")).
			Render(fmt.Sprintf("╭%s╮", strings.Repeat("─", m.platformVP.Width-2))),
	}

	layerCount := len(m.tag.Platforms[m.activeTabIndex].Layers)

	for i, layer := range m.tag.Platforms[m.activeTabIndex].Layers {
		size := humanReadableSize(layer.Size)

		if size != fmt.Sprintf("%dB", layer.Size) {
			size += fmt.Sprintf(" (%dB)", layer.Size)
		}

		out := []string{
			drawPercent(sizeBytes, layer.Size, m.platformVP.Width-3),
			"",
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().
					Faint(true).
					Width(9).
					Render("Digest:"),
				lipgloss.NewStyle().
					Width(m.platformVP.Width-13).
					Render(layer.Digest),
			),
			"",
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().
					Faint(true).
					Width(9).
					Render("Size:"),
				lipgloss.NewStyle().
					Width(m.platformVP.Width-13).
					Render(size),
			),
		}

		layers = append(layers,
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(lipgloss.Color("#383838")).
				BorderTop(false).
				BorderBottom(false).
				Render(lipgloss.JoinVertical(lipgloss.Top, out...)),
		)

		if i != layerCount-1 {
			layers = append(layers, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#383838")).
				Render(fmt.Sprintf("├%s┤", strings.Repeat("─", m.platformVP.Width-2))))
		}
	}

	if layerCount > 0 {
		layers = append(layers,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#383838")).
				Render(fmt.Sprintf("╰%s╯", strings.Repeat("─", m.platformVP.Width-2))),
		)
	}

	strs = append(strs, layers...)

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}

func drawPercent(total, current int, width int) string {
	if total <= 0 {
		return strings.Repeat("░", width-7) + "   0.0%"
	}

	rawPercent := (float64(current) / float64(total)) * 100

	displayPercent := math.Floor(rawPercent*10) / 10

	if current < total && displayPercent >= 100.0 {
		displayPercent = 99.9
	}

	if current >= total {
		displayPercent = 100.0
	}

	barWidth := max(width-7, 0)

	filledLength := int(math.Floor(float64(barWidth) * (displayPercent / 100)))
	emptyLength := barWidth - filledLength

	bar := strings.Repeat("█", filledLength) + strings.Repeat("░", emptyLength)

	return fmt.Sprintf("%s %6.1f%%", bar, displayPercent)
}
