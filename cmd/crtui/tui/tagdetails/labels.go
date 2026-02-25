package tagdetails

import (
	"maps"
	"slices"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/samber/lo"
)

func (m *Model) drawLabels(activePlatform models.Platform, width int) string {
	labelMap := maps.Clone(activePlatform.Config.Labels)

	if len(labelMap) == 0 {
		return ""
	}

	var (
		row1 [][2]any
		row2 [][2]any
		row3 [][2]any
		row4 [][2]any
		row5 [][2]any // other fields
	)

	if val, ok := labelMap["org.opencontainers.image.title"]; ok {
		row1 = append(row1, [2]any{"Title", val})
		delete(labelMap, "org.opencontainers.image.title")
	}

	if val, ok := labelMap["org.opencontainers.image.version"]; ok {
		row1 = append(row1, [2]any{"Version", val})
		delete(labelMap, "org.opencontainers.image.version")
	}

	if val, ok := labelMap["org.opencontainers.image.description"]; ok {
		row2 = append(row2, [2]any{"Description", val})
		delete(labelMap, "org.opencontainers.image.description")
	}

	if val, ok := labelMap["org.opencontainers.image.source"]; ok {
		row3 = append(row3, [2]any{"Source", val})
		delete(labelMap, "org.opencontainers.image.source")
	}

	if val, ok := labelMap["org.opencontainers.image.authors"]; ok {
		row4 = append(row4, [2]any{"Authors", val})
		delete(labelMap, "org.opencontainers.image.authors")
	}

	if val, ok := labelMap["org.opencontainers.image.created"]; ok {
		row4 = append(row4, [2]any{"Created", val})
		delete(labelMap, "org.opencontainers.image.created")
	}

	if val, ok := labelMap["org.opencontainers.image.licenses"]; ok {
		row4 = append(row4, [2]any{"Licenses", val})
		delete(labelMap, "org.opencontainers.image.licenses")
	}

	rows := make([]string, 0)

	if len(row1) > 0 {
		rows = append(rows, m.drawGrid(row1, 20, width))
	}

	if len(row2) > 0 {
		rows = append(rows, m.drawGrid(row2, 20, width))
	}

	if len(row3) > 0 {
		rows = append(rows, m.drawGrid(row3, 20, width))
	}

	if len(row4) > 0 {
		rows = append(rows, m.drawGrid(row4, 20, width))
	}

	sortedLabelKeys := lo.MapToSlice(labelMap, func(k, _ string) string { return k })
	slices.Sort(sortedLabelKeys)

	for _, k := range sortedLabelKeys {
		row5 = append(row5, [2]any{k, labelMap[k]})
	}

	if len(row5) > 0 {
		rows = append(rows, m.drawGrid(row5, width, width))
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Render("❯ LABELS"),
	}

	strs = append(strs, rows...)

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
