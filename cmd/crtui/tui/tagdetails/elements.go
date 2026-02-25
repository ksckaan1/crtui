package tagdetails

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func (m *Model) drawGrid(fields [][2]any, minElementWidth, rowWidth int) string {
	n := len(fields)
	if n == 0 {
		return ""
	}

	maxCn := max(rowWidth/minElementWidth, 1)

	cn := min(n, maxCn) // column count

	cw := rowWidth / cn // column width

	cellStyle := lipgloss.NewStyle().
		Width(cw).
		Align(lipgloss.Left).
		MarginBottom(1)

	var rows []string
	var currentRow []string

	for i, field := range fields {
		renderedCell := cellStyle.Render(m.drawField(field[0], field[1], cw))
		currentRow = append(currentRow, renderedCell)

		if (i+1)%cn == 0 || i+1 == n {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *Model) drawField(title, content any, width int) string {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false).
		BorderLeft(true).
		PaddingLeft(1).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().
				Faint(true).
				Width(width-2).
				Render(fmt.Sprintf("%v", title)),
			lipgloss.NewStyle().
				Width(width-2).
				Render(fmt.Sprintf("%v", content)),
		))
}
