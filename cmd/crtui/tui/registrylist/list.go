package registrylist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

type Registry struct {
	URL      string
	Username string
	Password string
	Status   string
}

func (i *Registry) Title() string       { return i.URL }
func (i *Registry) Description() string { return "@" + i.Username }
func (i *Registry) FilterValue() string { return i.URL + " " + i.Username }

type registryListDelegate struct {
	getSpinnerFunc func() string
}

func (d *registryListDelegate) Height() int                             { return 3 }
func (d *registryListDelegate) Spacing() int                            { return 0 }
func (d *registryListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d *registryListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	r, ok := item.(*Registry)
	if !ok {
		return
	}

	isSelected := m.Index() == index

	status := offlineDot

	switch r.Status {
	case "online":
		status = onlineDot

	case "loading":
		status = d.getSpinnerFunc()
		if isSelected {
			status = lipgloss.NewStyle().
				Foreground(primaryColor).
				Render(status)
		}

	case "unauth":
		status = unauthDot
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EAEAEA")).Faint(true)

	if isSelected {
		titleStyle = titleStyle.Foreground(primaryColor)
		subtitleStyle = subtitleStyle.Foreground(primaryColor)
	}

	uri := r.URL

	uri = strings.TrimPrefix(uri, "https://")
	uri = strings.TrimPrefix(uri, "http://")

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		" ",
		status,
		lo.Ternary(r.Status == "loading", "", " "),
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(uri),
			subtitleStyle.Render("@"+r.Username),
		),
	)

	leftBorder := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, false, false, true)

	if isSelected {
		leftBorder = leftBorder.Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderLeftForeground(primaryColor)
	}

	out = leftBorder.Render(out)

	out = lipgloss.NewStyle().MarginBottom(1).Render(out)

	fmt.Fprint(w, out)
}
