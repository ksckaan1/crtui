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
	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/samber/lo"
)

var _ tea.Model = (*Model)(nil)

type Model struct {
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

func NewModel(registries []*Registry) *Model {
	items := lo.Map(registries, func(item *Registry, _ int) list.Item {
		return item
	})

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.SetShowStatusBar(false)

	rlw := ui.NewWindow()
	rlw.SetLeftTitle("Container Registries")
	rlw.SetRightTitle(ui.CountKind(len(items), "registry", "registries"))

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
		keysWindow:             kw,
		statusesWindow:         sw,
		registryListWindow:     rlw,
		list:                   l,
		spinner:                sp,
		registries:             registries,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),
		status:                 ui.NewStatus(),
	}
}

func (m *Model) Init() tea.Cmd {
	m.list.SetDelegate(&registryListDelegate{
		getSpinnerFunc: func() string {
			return m.spinner.View()
		},
	})

	cmds := []tea.Cmd{m.spinner.Tick, tea.RequestWindowSize}
	cmds = append(cmds, m.fetchRegistries()...)

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.SelectRegistry) &&
			m.list.FilterState() != list.Filtering &&
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

		case key.Matches(msg, keyMap.Refresh) && m.list.FilterState() != list.Filtering:
			return m, tea.Batch(m.fetchRegistries()...)

		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.keysWindow.SetWidth(m.width - 44)
		m.registryListWindow.SetWidth(m.width - 2)
		m.registryListWindow.SetHeight(m.height - 8 - m.status.Height())

	case registryResult:
		m.registries[msg.index].Status = msg.status
		m.registries[msg.index].SupportsHTTP3 = msg.supportsHTTP3
	}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)

	return m, tea.Batch(spinnerCmd, listCmd)
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

	m.registryListWindow.SetContent(func(width, height int) string {
		m.list.SetWidth(width)
		m.list.SetHeight(height)

		return m.list.View()
	})

	out := lipgloss.NewStyle().
		Padding(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.drawHeader(),
				m.registryListWindow.View(),
				m.status.View(),
			),
		)

	fmt.Fprint(sb, out)

	v.SetContent(sb.String())

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
