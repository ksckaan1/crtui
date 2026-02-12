package registrylist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/samber/lo"
)

type Registry struct {
	URL           string
	Username      string
	Password      string
	Status        registrystatus.RegistryStatus
	SupportsHTTP3 bool
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
	case registrystatus.Online:
		status = onlineDot

	case registrystatus.Loading:
		status = d.getSpinnerFunc()
		if isSelected {
			status = lipgloss.NewStyle().
				Foreground(ui.PrimaryColor).
				Render(status)
		}

	case registrystatus.Unauth:
		status = unauthDot
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).Faint(true)
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EAEAEA")).Faint(true)

	if isSelected {
		titleStyle = titleStyle.Foreground(lipgloss.Color("#FFFFFF")).Faint(false).Bold(true)
		subtitleStyle = subtitleStyle.Foreground(lipgloss.Color("#FFFFFF")).Faint(false).Bold(true)
	}

	uri := strings.TrimPrefix(r.URL, "https://")
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
			BorderLeftForeground(lipgloss.Color("#FFFFFF"))
	}

	out = leftBorder.Render(out)

	out = lipgloss.NewStyle().MarginBottom(1).Render(out)

	fmt.Fprint(w, out)
}
