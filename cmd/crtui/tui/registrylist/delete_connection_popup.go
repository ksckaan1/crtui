package registrylist

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/config"
)

var _ tea.Model = (*DeleteConnectionPopup)(nil)

type DeleteConnectionPopup struct {
	cfg                    *config.Config
	status                 *ui.Status
	back                   background
	backgroundText         string
	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	ov                     *ui.Overlay
	registry               *Registry

	// ui elements
	popup                *ui.Window
	deleteButton         *ui.Button
	cancelButton         *ui.Button
	focusables           []focusable
	activeFocusableIndex int
}

func NewDeleteConnectionPopup(
	cfg *config.Config,
	back background,
	status *ui.Status,
	registry *Registry,
) *DeleteConnectionPopup {
	popup := ui.NewWindow()
	popup.SetWidth(40)
	popup.SetHeight(10)
	popup.SetLeftTitle("Delete Connection")
	popup.SetRightTitle("esc")
	popup.SetBorderColor(ui.PrimaryColor)
	popup.SetLeftTitleColor(lipgloss.White)
	popup.SetRightTitleColor(lipgloss.White)

	cancelButton := ui.NewButton(" Cancel ")
	cancelButton.SetNormalColors(lipgloss.Color("#FFFFFF"), nil)

	deleteButton := ui.NewButton(" Delete ")
	deleteButton.SetNormalColors(lipgloss.Red, nil)
	deleteButton.SetWidth(27)

	return &DeleteConnectionPopup{
		cfg:                    cfg,
		status:                 status,
		back:                   back,
		backgroundText:         back.View().Content,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),
		ov:                     ui.NewOverlay(status),
		registry:               registry,

		popup:        popup,
		deleteButton: &deleteButton,
		cancelButton: &cancelButton,
		focusables: []focusable{
			&cancelButton,
			&deleteButton,
		},
	}
}

func (m *DeleteConnectionPopup) Init() tea.Cmd {
	return tea.Batch(
		m.cancelButton.Focus(),
		tea.RequestWindowSize,
	)
}

func (m *DeleteConnectionPopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.back.UpdateSize(msg)
		m.backgroundText = m.back.View().Content

	case tea.KeyMsg:
		switch {
		case msg.String() == "esc":
			for i := range m.focusables {
				m.focusables[i].Blur()
			}
			m.status.SetStatus(ui.Empty, "")
			return nav.Navigate(m.back)

		case msg.String() == "ctrl+c":
			return m, tea.Quit

		case msg.String() == "left":
			for i := range m.focusables {
				m.focusables[i].Blur()
			}
			m.activeFocusableIndex--
			if m.activeFocusableIndex < 0 {
				m.activeFocusableIndex = 0
			}
			cmds = append(cmds, m.focusables[m.activeFocusableIndex].Focus())

		case msg.String() == "right":
			for i := range m.focusables {
				m.focusables[i].Blur()
			}
			m.activeFocusableIndex++
			if m.activeFocusableIndex > len(m.focusables)-1 {
				m.activeFocusableIndex = len(m.focusables) - 1
			}
			cmds = append(cmds, m.focusables[m.activeFocusableIndex].Focus())

		case msg.String() == "enter" &&
			m.activeFocusableIndex == 0:
			return nav.Navigate(m.back)

		case msg.String() == "enter" &&
			m.activeFocusableIndex == 1:
			return m.onDelete()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *DeleteConnectionPopup) onDelete() (tea.Model, tea.Cmd) {
	err := m.cfg.DeleteAuth(m.registry.URL, m.registry.Username)
	if err != nil {
		return m, m.status.SetStatus(ui.Error, err.Error())
	}

	m.status.SetStatus(ui.Info, "Connection deleted")
	return nav.Navigate(m.back)
}

func (m *DeleteConnectionPopup) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	m.ov.SetBack(m.backgroundText)

	username := "anonymous"

	if m.registry.Username != "" {
		username = fmt.Sprintf("@%s", m.registry.Username)
	}

	m.popup.SetContent(func(width, height int) string {

		popupContent := lipgloss.NewStyle().
			Padding(1).
			Width(width).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Render("Are you sure you want to delete connection?"),
					"",
					lipgloss.NewStyle().
						Bold(true).
						Render(ui.TrimURLScheme(m.registry.URL)),
					lipgloss.NewStyle().
						Faint(true).
						Render(username),
					"",
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						m.cancelButton.View(), " ", m.deleteButton.View(),
					),
				),
			)

		return popupContent
	})

	m.ov.SetFront(m.popup.View())

	v.SetContent(m.ov.View())

	return v
}
