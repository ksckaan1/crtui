package tagdetails

import (
	"maps"
	"slices"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/samber/lo"
)

func (m *TagDetailsScreenModel) drawLabels(activePlatform models.Platform, width int) string {
	labelMap := maps.Clone(activePlatform.Config.Labels)

	if len(labelMap) == 0 {
		return ""
	}

	strs := []string{
		lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Bold(true).
			MarginBottom(1).
			Render("❯ LABELS"),
	}

	sortedLabelKeys := lo.MapToSlice(labelMap, func(k, _ string) string { return k })
	slices.Sort(sortedLabelKeys)

	for _, k := range sortedLabelKeys {
		strs = append(strs, m.drawField(k, labelMap[k], width), "")
	}

	return lipgloss.JoinVertical(lipgloss.Top, strs...)
}
