package registrydetails

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

type Repository struct {
	Name string
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

	if d.selectedRepository != nil && *d.selectedRepository == r.Name {
		name = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(m.Width()-3).Render(name),
			lipgloss.NewStyle().Foreground(ui.PrimaryColor).Render("→"),
		)
	}

	fmt.Fprint(w, itemStyle.Render(name))
}
