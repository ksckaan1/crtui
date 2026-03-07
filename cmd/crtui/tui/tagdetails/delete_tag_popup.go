package tagdetails

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/config"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
)

var _ tea.Model = (*DeleteTagPopupModel)(nil)

type DeleteTagPopupModel struct {
	rc                     *registryclient.RegistryClient
	cfg                    *config.Config
	status                 *ui.Status
	back                   *TagDetailsScreenModel
	backgroundText         string
	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	ov                     *ui.Overlay
	repositoryName         string
	tagName                string
	isLoading              bool

	// ui elements
	popup                *ui.Window
	deleteButton         *ui.Button
	cancelButton         *ui.Button
	focusables           []focusable
	activeFocusableIndex int
}

type focusable interface {
	Focus() tea.Cmd
	Blur()
	Reset()
}

func NewDeleteTagPopupModel(
	rc *registryclient.RegistryClient,
	back *TagDetailsScreenModel,
	status *ui.Status,
	repositoryName string,
	tagName string,
) *DeleteTagPopupModel {

	popup := ui.NewWindow()
	popup.SetWidth(40)
	popup.SetHeight(8)
	popup.SetLeftTitle("Delete Tag")
	popup.SetRightTitle("esc")
	popup.SetBorderColor(ui.PrimaryColor)
	popup.SetLeftTitleColor(lipgloss.White)
	popup.SetRightTitleColor(lipgloss.White)

	cancelButton := ui.NewButton(" Cancel ")
	cancelButton.SetNormalColors(lipgloss.Color("#FFFFFF"), nil)

	deleteButton := ui.NewButton(" Delete ")
	deleteButton.SetNormalColors(lipgloss.Red, nil)
	deleteButton.SetWidth(27)

	return &DeleteTagPopupModel{
		rc:                     rc,
		status:                 status,
		back:                   back,
		backgroundText:         back.View().Content,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(60, 24),
		ov:                     ui.NewOverlay(status),
		repositoryName:         repositoryName,
		tagName:                tagName,

		popup:        popup,
		deleteButton: &deleteButton,
		cancelButton: &cancelButton,
		focusables: []focusable{
			&cancelButton,
			&deleteButton,
		},
	}
}

func (m *DeleteTagPopupModel) Init() tea.Cmd {
	return tea.Batch(
		m.cancelButton.Focus(),
		tea.RequestWindowSize,
	)
}

func (m *DeleteTagPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			m.isLoading = true

			return m, tea.Batch(
				m.status.SetStatus(ui.Info, "Deleting tag..."),
				m.deleteTag(),
			)
		}

	case deleteTagMsg:
		m.isLoading = false

		if msg.err != nil {
			return m, m.status.SetStatus(ui.Error, msg.err.Error())
		}

		m.status.SetStatus(ui.OK, "Tag deleted successfully")
		return nav.Navigate(m.back.back)
	}

	return m, tea.Batch(cmds...)
}

func (m *DeleteTagPopupModel) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	m.ov.SetBack(m.backgroundText)

	m.popup.SetContent(func(width, height int) string {
		popupContent := lipgloss.NewStyle().
			Padding(1).
			Width(width).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Render("Are you sure you want to delete tag?"),
					"",
					lipgloss.NewStyle().
						Bold(true).
						Render(fmt.Sprintf("%s/%s", m.repositoryName, m.tagName)),
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
