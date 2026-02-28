package tagdetails

import (
	"charm.land/lipgloss/v2"
	"github.com/samber/lo"
)

func (m *TagDetailsScreenModel) drawField(title, content string, width int) string {
	keyField := lipgloss.NewStyle().
		Faint(true).
		Width(width - 2)

	valueField := lipgloss.NewStyle().
		Width(width - 2)

	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false).
		BorderLeft(true).
		PaddingLeft(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				keyField.Render(title),
				lo.Ternary(
					content != "",
					valueField.Render(content),
					keyField.Render("<empty string>"),
				),
			),
		)
}
