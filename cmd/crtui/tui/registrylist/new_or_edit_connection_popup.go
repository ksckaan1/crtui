package registrylist

import (
	"context"
	"errors"
	"net/url"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/config"
	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
)

var _ tea.Model = (*NewOrEditConnectionPopup)(nil)

type NewOrEditConnectionPopup struct {
	cfg                    *config.Config
	status                 *ui.Status
	back                   background
	backgroundText         string
	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	ov                     *ui.Overlay

	// ui elements
	newConnectionPopup   *ui.Window
	crURLtextInput       *textinput.Model
	isAuthRequired       *ui.Checkbox
	usernameTextInput    *textinput.Model
	passwordTextInput    *textinput.Model
	testButton           *ui.Button
	connectButton        *ui.Button
	focusables           []focusable
	activeFocusableIndex int
}

type background interface {
	tea.Model
	UpdateSize(tea.WindowSizeMsg)
}

func NewNewOrEditConnectionPopup(
	cfg *config.Config,
	back background,
	status *ui.Status,
	registry *Registry,
) *NewOrEditConnectionPopup {
	ncw := ui.NewWindow()
	ncw.SetWidth(40)
	ncw.SetHeight(17)
	ncw.SetLeftTitle("New Connection")
	ncw.SetRightTitle("esc")
	ncw.SetBorderColor(ui.PrimaryColor)
	ncw.SetLeftTitleColor(lipgloss.White)
	ncw.SetRightTitleColor(lipgloss.White)

	crURLti := textinput.New()
	crURLti.Placeholder = "https://example.com"
	crURLti.Validate = func(s string) error {
		u, err := url.Parse(s)
		if err != nil {
			return err
		}
		if u.Scheme != "https" {
			return errors.New("only https connections are valid")
		}
		return err
	}

	isAuthRequired := ui.NewCheckbox()

	usernameti := textinput.New()
	usernameti.Placeholder = "username"
	usernameti.Prompt = "Insert:  "

	passwordti := textinput.New()
	passwordti.Placeholder = "password"
	passwordti.Prompt = "Insert:  "
	passwordti.EchoMode = textinput.EchoPassword

	connectBtn := ui.NewButton("Create")
	testBtn := ui.NewButton("Test")
	testBtn.SetNormalColors(lipgloss.White, lipgloss.Color("#444444"))

	if registry != nil {
		ncw.SetLeftTitle("Edit Connection")
		crURLti.SetValue(registry.URL)
		usernameti.SetValue(registry.Username)
		passwordti.SetValue(registry.Password)
		if registry.Username != "" || registry.Password != "" {
			isAuthRequired.SetValue(true)
		}
	}

	return &NewOrEditConnectionPopup{
		cfg:                    cfg,
		status:                 status,
		back:                   back,
		backgroundText:         back.View().Content,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),
		ov:                     ui.NewOverlay(status),

		newConnectionPopup: ncw,
		crURLtextInput:     &crURLti,
		isAuthRequired:     &isAuthRequired,
		usernameTextInput:  &usernameti,
		passwordTextInput:  &passwordti,
		connectButton:      &connectBtn,
		testButton:         &testBtn,
		focusables: []focusable{
			&crURLti,
			&isAuthRequired,
			&usernameti,
			&passwordti,
			&testBtn,
			&connectBtn,
		},
	}
}

func (m *NewOrEditConnectionPopup) Init() tea.Cmd {
	return tea.Batch(
		m.crURLtextInput.Focus(),
		tea.RequestWindowSize,
	)
}

func (m *NewOrEditConnectionPopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case msg.String() == "up":
			for i := range m.focusables {
				m.focusables[i].Blur()
			}
			m.activeFocusableIndex--
			if m.activeFocusableIndex < 0 {
				m.activeFocusableIndex = len(m.focusables) - 1
			}
			if !m.isAuthRequired.Value() && (m.activeFocusableIndex == 2 || m.activeFocusableIndex == 3) {
				m.activeFocusableIndex = 1
			}
			cmds = append(cmds, m.focusables[m.activeFocusableIndex].Focus())

		case msg.String() == "down" || msg.String() == "tab":
			for i := range m.focusables {
				m.focusables[i].Blur()
			}
			m.activeFocusableIndex++
			if !m.isAuthRequired.Value() && (m.activeFocusableIndex == 2 || m.activeFocusableIndex == 3) {
				m.activeFocusableIndex = 4
			}
			if m.activeFocusableIndex > len(m.focusables)-1 {
				m.activeFocusableIndex = 0
			}
			cmds = append(cmds, m.focusables[m.activeFocusableIndex].Focus())

		case msg.String() == "enter" &&
			m.activeFocusableIndex == 4:
			return m.onTest()

		case msg.String() == "enter" &&
			m.activeFocusableIndex == 5:
			return m.onCreate()
		}
	}

	crURLtextInput, crURLtextInputCmd := m.crURLtextInput.Update(msg)
	*m.crURLtextInput = crURLtextInput
	isAuthRequired, isAuthRequiredCmd := m.isAuthRequired.Update(msg)
	*m.isAuthRequired = isAuthRequired
	usernameTextInput, usernameTextInputCmd := m.usernameTextInput.Update(msg)
	*m.usernameTextInput = usernameTextInput
	passwordTextInput, passwordTextInputCmd := m.passwordTextInput.Update(msg)
	*m.passwordTextInput = passwordTextInput
	cmds = append(cmds,
		crURLtextInputCmd,
		isAuthRequiredCmd,
		usernameTextInputCmd,
		passwordTextInputCmd,
	)

	return m, tea.Batch(cmds...)
}

func (m *NewOrEditConnectionPopup) onTest() (tea.Model, tea.Cmd) {
	if m.crURLtextInput.Err != nil || m.crURLtextInput.Value() == "" {
		return m, m.status.SetStatus(ui.Error, "Insert valid container registry url")
	}

	username := m.usernameTextInput.Value()
	password := m.passwordTextInput.Value()

	if m.isAuthRequired.Value() &&
		(username == "" || password == "") {
		return m, m.status.SetStatus(ui.Error, "Insert valid username and password")
	}

	if !m.isAuthRequired.Value() {
		username = ""
		password = ""
	}

	rc := registryclient.New(
		m.crURLtextInput.Value(),
		username,
		password,
	)

	ri, err := rc.GetRegistryInfo(context.Background())
	if err != nil {
		return m, m.status.SetStatus(ui.Error, err.Error())
	}

	if ri.Status == registrystatus.Unauth {
		return m, m.status.SetStatus(ui.Error, "Authentication failed")
	}

	if ri.Status == registrystatus.Invalid {
		return m, m.status.SetStatus(ui.Error, "Invalid container registry")
	}

	return m, m.status.SetStatus(ui.Info, "Connection successful")
}

func (m *NewOrEditConnectionPopup) onCreate() (tea.Model, tea.Cmd) {
	if m.crURLtextInput.Err != nil || m.crURLtextInput.Value() == "" {
		return m, m.status.SetStatus(ui.Error, "Insert valid container registry url")
	}

	username := m.usernameTextInput.Value()
	password := m.passwordTextInput.Value()

	if m.isAuthRequired.Value() &&
		(username == "" || password == "") {
		return m, m.status.SetStatus(ui.Error, "Insert valid username and password")
	}

	if !m.isAuthRequired.Value() {
		username = ""
		password = ""
	}

	err := m.cfg.SetAuth(
		m.crURLtextInput.Value(),
		username,
		password,
	)
	if err != nil {
		return m, m.status.SetStatus(ui.Error, err.Error())
	}

	cmds := []tea.Cmd{m.status.SetStatus(ui.Info, "Connection created")}

	model, cmd := nav.Navigate(m.back)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (m *NewOrEditConnectionPopup) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	m.ov.SetBack(m.backgroundText)

	m.newConnectionPopup.SetContent(func(width, height int) string {
		innerWidth := width - 3

		m.crURLtextInput.SetWidth(innerWidth - 9)
		m.usernameTextInput.SetWidth(innerWidth - 9)
		m.passwordTextInput.SetWidth(innerWidth - 9)
		m.connectButton.SetWidth(innerWidth)
		m.testButton.SetWidth(innerWidth)

		if m.crURLtextInput.Value() == "" {
			m.crURLtextInput.Prompt = lipgloss.NewStyle().
				Foreground(lipgloss.White).
				Render("Insert:  ")
		} else if m.crURLtextInput.Err != nil {
			m.crURLtextInput.Prompt = lipgloss.NewStyle().
				Foreground(lipgloss.Red).
				Render("Invalid: ")
		} else {
			m.crURLtextInput.Prompt = lipgloss.NewStyle().
				Foreground(lipgloss.Green).
				Render("Valid:   ")
		}

		authInputs := lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			lipgloss.NewStyle().
				Faint(true).
				Render("USERNAME"),
			m.usernameTextInput.View(),
			"",
			lipgloss.NewStyle().
				Faint(true).
				Render("PASSWORD"),
			m.passwordTextInput.View(),
		)

		if !m.isAuthRequired.Value() {
			authInputs = lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.Color("#444444")).
				Render(ansi.Strip(authInputs))
		}

		popupContent := lipgloss.NewStyle().
			Padding(1).
			Width(width).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Render("CONTAINER REGISTRY URL"),
					m.crURLtextInput.View(),
					"",
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						lipgloss.NewStyle().
							Faint(true).
							Render("AUTH:"),
						" ",
						m.isAuthRequired.View(),
					),
					authInputs,
					"",
					m.testButton.View(),
					"",
					m.connectButton.View(),
				),
			)

		return popupContent
	})

	m.ov.SetFront(m.newConnectionPopup.View())

	v.SetContent(m.ov.View())

	return v
}
