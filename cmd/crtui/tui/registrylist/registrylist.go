package registrylist

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/figlet"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrydetails"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/config"
	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/samber/lo"
)

type focusable interface {
	Focus() tea.Cmd
	Blur()
	Reset()
}

var _ tea.Model = (*Model)(nil)

type Model struct {
	cfg                *config.Config
	keysWindow         *ui.KeysWindow
	statusesWindow     *ui.Window
	registryListWindow *ui.Window
	list               list.Model
	spinner            spinner.Model

	registries []*Registry
	cursor     int

	width  int
	height int

	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	status                 *ui.Status
}

func NewModel(cfg *config.Config) *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.SetShowStatusBar(false)

	rlw := ui.NewWindow()

	kw := ui.NewKeysWindow()
	kw.SetHeight(5)
	kw.SetKeyMap(keyMap)

	sw := ui.NewWindow()
	sw.SetWidth(14)
	sw.SetHeight(5)

	sw.SetContent(func(width, height int) string {
		statusContainer := lipgloss.NewStyle().
			PaddingLeft(1)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

		onlineStatus := fmt.Sprintf("%s %s", onlineDot, descStyle.Render("Online"))
		offlineStatus := fmt.Sprintf("%s %s", offlineDot, descStyle.Render("Offline"))
		unauthStatus := fmt.Sprintf("%s %s", unauthDot, descStyle.Render("Unauth"))

		statusContent := lipgloss.JoinVertical(
			lipgloss.Left,
			onlineStatus,
			offlineStatus,
			unauthStatus,
		)

		return statusContainer.Render(statusContent)
	})

	return &Model{
		cfg:                    cfg,
		keysWindow:             kw,
		statusesWindow:         sw,
		registryListWindow:     rlw,
		list:                   l,
		spinner:                sp,
		registries:             nil,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),
		status:                 ui.NewStatus(),
	}
}

func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick, tea.RequestWindowSize}

	m.list.SetDelegate(&registryListDelegate{
		getSpinnerFunc: func() string {
			return m.spinner.View()
		},
	})

	auths, err := m.cfg.ListAuths()
	if err != nil {
		cmds = append(cmds,
			m.status.SetStatus(ui.Error, err.Error()),
		)
		return tea.Batch(cmds...)
	}

	m.registries = lo.Map(auths, func(item *config.Auth, _ int) *Registry {
		return &Registry{
			URL:          item.URL,
			Username:     item.Username,
			Password:     item.Password,
			Status:       registrystatus.Loading,
			AutoDetected: item.AutoDetected,
		}
	})

	m.registryListWindow.SetRightTitle(
		ui.CountKind(len(m.registries),
			"registry",
			"registries",
		),
	)

	m.list.SetItems(lo.Map(m.registries, func(item *Registry, _ int) list.Item {
		return item
	}))

	cmds = append(cmds, m.fetchRegistries()...)

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// SELECT REGISTRY
		case key.Matches(msg, keyMap.SelectRegistry) &&
			!m.isTyping() &&
			m.registries[m.list.Index()].Status == registrystatus.Online:
			reg := m.registries[m.list.Index()]
			m.status.SetStatus(ui.Empty, "")
			return nav.Navigate(
				registrydetails.NewRegistryDetails(
					&registrydetails.Registry{
						URL:            reg.URL,
						Username:       reg.Username,
						Password:       reg.Password,
						SupportsHTTPS3: reg.SupportsHTTP3,
					},
					m,
					m.status,
				))

		// REFRESH REGISTRY LIST
		case key.Matches(msg, keyMap.Refresh) &&
			!m.isTyping():
			return m, tea.Batch(m.fetchRegistries()...)

		case key.Matches(msg, keyMap.New) &&
			!m.isTyping():
			return nav.Navigate(NewNewConnectionPopup(m.cfg, m, m.status))

		// QUIT PROGRAM
		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.UpdateSize(msg)

	case registryResult:
		m.registries[msg.index].Status = msg.status
		m.registries[msg.index].SupportsHTTP3 = msg.supportsHTTP3
	}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	cmds = append(cmds, spinnerCmd)

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) UpdateSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
	m.keysWindow.SetWidth(m.width - 44)
	m.registryListWindow.SetWidth(m.width - 2)
	m.registryListWindow.SetHeight(m.height - 8 - m.status.Height())
}

func (m *Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	sb := &strings.Builder{}

	leftTitle := "Container Registries"

	if m.list.IsFiltered() {
		leftTitle += lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Render(" /" + m.list.FilterValue())
	}

	m.registryListWindow.SetLeftTitle(leftTitle)

	m.registryListWindow.SetContent(func(width, height int) string {
		m.list.SetWidth(width)
		m.list.SetHeight(height)

		return m.list.View()
	})

	out := lipgloss.NewStyle().
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.drawHeader(),
				m.registryListWindow.View(),
			),
		)

	fmt.Fprint(sb, out)

	view := sb.String()

	v.SetContent(
		lipgloss.NewStyle().
			Padding(1).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					view,
					m.status.View(),
				),
			),
	)

	return v
}

func (m *Model) drawHeader() string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		m.keysWindow.View(),
		" ",
		m.statusesWindow.View(),
	)

	return out
}

func (m *Model) fetchRegistries() []tea.Cmd {
	return lo.Map(m.registries, func(item *Registry, i int) tea.Cmd {
		item.Status = registrystatus.Loading
		return fetchRegistry(i, item)
	})
}

func (m *Model) isTyping() bool {
	return m.list.FilterState() == list.Filtering
}
