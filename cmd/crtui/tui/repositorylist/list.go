package repositorylist

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

type Repository struct {
	Name       string
	Visibility string
}

func (i *Repository) Title() string       { return i.Name }
func (i *Repository) Description() string { return i.Name }
func (i *Repository) FilterValue() string { return i.Name }

type repositoryListDelegate struct {
	selectedRepository *string
}

func (d *repositoryListDelegate) Height() int                             { return 2 }
func (d *repositoryListDelegate) Spacing() int                            { return 0 }
func (d *repositoryListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d *repositoryListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	r, ok := item.(*Repository)
	if !ok {
		return
	}

	isSelected := m.Index() == index

	itemStyle := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, false, false, true).
		Faint(true).
		MarginBottom(1).
		PaddingLeft(1)

	if isSelected {
		itemStyle = itemStyle.
			Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderLeftForeground(ui.PrimaryColor).
			Faint(false)
	}

	name := r.Name

	visibilityBadge := ""

	if r.Visibility != "" {
		backgroundColor := lipgloss.Color("#3FB950")

		if r.Visibility == "private" {
			backgroundColor = lipgloss.Color("#F85149")
		}

		visibilityBadge = lipgloss.NewStyle().
			Padding(0, 1).
			Background(backgroundColor).
			Foreground(lipgloss.White).
			Bold(true).
			Render(strings.ToUpper(r.Visibility))
	}

	width := m.Width() - lipgloss.Width(visibilityBadge) - 4
	if width < 0 {
		width = 0
	}

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(width).Render(name),
		visibilityBadge,
	)

	if d.selectedRepository != nil && *d.selectedRepository == r.Name {
		line = lipgloss.JoinHorizontal(
			lipgloss.Top,
			line,
			lipgloss.NewStyle().Foreground(ui.PrimaryColor).Render("→"),
		)
	}

	fmt.Fprint(w, itemStyle.Render(line))
}
