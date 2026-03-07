package taglist

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
	"github.com/samber/lo"
)

var _ tea.Model = (*DeleteTagsPopup)(nil)

type DeleteTagsPopup struct {
	rc                     *registryclient.RegistryClient
	status                 *ui.Status
	back                   *TagListScreenModel
	backgroundText         string
	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	ov                     *ui.Overlay
	repositoryName         string
	tagNames               []string

	// ui elements
	spinner              spinner.Model
	popup                *ui.Window
	deleteButton         *ui.Button
	cancelButton         *ui.Button
	focusables           []focusable
	activeFocusableIndex int
	tagNamesVP           viewport.Model
}

type focusable interface {
	Focus() tea.Cmd
	Blur()
	Reset()
}

func NewDeleteTagsPopup(
	rc *registryclient.RegistryClient,
	back *TagListScreenModel,
	status *ui.Status,
	repositoryName string,
	tagNames []string,
	spinner spinner.Model,
) *DeleteTagsPopup {
	leftTitle := "Delete Tag"

	if len(tagNames) > 1 {
		leftTitle = "Delete Tags"
	}

	if len(tagNames) == 0 {
		leftTitle = "Delete Repository"
	}

	popup := ui.NewWindow()
	popup.SetWidth(40)
	popup.SetHeight(9)
	popup.SetLeftTitle(leftTitle)
	popup.SetRightTitle("esc")
	popup.SetBorderColor(ui.PrimaryColor)
	popup.SetLeftTitleColor(lipgloss.White)
	popup.SetRightTitleColor(lipgloss.White)

	cancelButton := ui.NewButton(" Cancel ")
	cancelButton.SetNormalColors(lipgloss.Color("#FFFFFF"), nil)

	deleteButton := ui.NewButton(" Delete ")
	deleteButton.SetNormalColors(lipgloss.Red, nil)
	deleteButton.SetWidth(27)

	return &DeleteTagsPopup{
		rc:                     rc,
		status:                 status,
		back:                   back,
		backgroundText:         back.View().Content,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(60, 24),
		ov:                     ui.NewOverlay(status),
		repositoryName:         repositoryName,
		tagNames:               tagNames,

		spinner:      spinner,
		popup:        popup,
		deleteButton: &deleteButton,
		cancelButton: &cancelButton,
		focusables: []focusable{
			&cancelButton,
			&deleteButton,
		},
	}
}

func (m *DeleteTagsPopup) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.cancelButton.Focus(),
		tea.RequestWindowSize,
	}

	if len(m.tagNames) == 0 {
		cmds = append(cmds, m.fetchTagList())
	}

	return tea.Batch(cmds...)
}

func (m *DeleteTagsPopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case tagListResult:
		if msg.err != nil {
			return m, m.status.SetStatus(ui.Error, msg.err.Error())
		}

		m.tagNames = msg.tagList

	case deleteTagsResult:
		if msg.err != nil {
			return m, m.status.SetStatus(ui.Error, msg.err.Error())
		}

		m.status.SetStatus(ui.OK, "Deleted selected tags")
		return nav.Navigate(m.back)
	}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	return m, tea.Batch(cmds...)
}

func (m *DeleteTagsPopup) onDelete() (tea.Model, tea.Cmd) {
	if len(m.tagNames) == 0 {
		return m, nil
	}

	return m, tea.Batch(
		m.status.SetStatus(ui.Info, "Deleting selected tags..."),
		m.deleteTags(),
	)
}

func (m *DeleteTagsPopup) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	m.ov.SetBack(m.backgroundText)

	tags := lo.Ternary(
		len(m.tagNames) == 0,
		m.spinner.View(),
		ui.CountKind(len(m.tagNames), "tag", "tags"),
	)

	m.popup.SetContent(func(width, height int) string {

		popupContent := lipgloss.NewStyle().
			Padding(1).
			Width(width).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Render("Are you sure you want to delete?"),
					"",
					m.repositoryName,
					tags,
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
