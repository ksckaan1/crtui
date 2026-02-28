package registrylist

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	AutoDetected  bool
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

	username := lo.Ternary(r.Username != "", "@"+r.Username, "anonymous")

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		" ",
		status,
		lo.Ternary(r.Status == "loading", "", " "),
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(ui.TrimURLScheme(r.URL)),
			subtitleStyle.Render(username),
		),
	)

	leftBorder := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, false, false, true)

	if isSelected {
		leftBorder = leftBorder.Border(lipgloss.ThickBorder(), false, false, false, true).
			BorderLeftForeground(ui.PrimaryColor)
	}

	out = leftBorder.Render(out)

	autoDetectedBadge := ""

	if r.AutoDetected {
		autoDetectedBadge = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.White).
			Foreground(lipgloss.Black).
			Bold(true).
			Render("AUTO-DETECTED")
	}

	out = lipgloss.NewStyle().
		MarginBottom(1).
		Width(m.Width() - lipgloss.Width(autoDetectedBadge) - 2).
		Render(out)

	fmt.Fprint(w,
		lipgloss.JoinHorizontal(lipgloss.Top,
			out,
			autoDetectedBadge,
		),
	)
}
