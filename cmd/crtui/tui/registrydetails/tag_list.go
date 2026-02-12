package registrydetails

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

type Tag struct {
	Name string
}

func (i *Tag) Title() string       { return i.Name }
func (i *Tag) Description() string { return i.Name }
func (i *Tag) FilterValue() string { return i.Name }

type tagListDelegate struct{}

func (d *tagListDelegate) Height() int                             { return 2 }
func (d *tagListDelegate) Spacing() int                            { return 0 }
func (d *tagListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d *tagListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	r, ok := item.(*Tag)
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

	fmt.Fprint(w, itemStyle.Render(r.Name))
}
