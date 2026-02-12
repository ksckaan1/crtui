package registrylist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/figlet"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrydetails"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/samber/lo"
)

var _ tea.Model = (*Model)(nil)

type listKeyMap struct {
	refresh        key.Binding
	selectRegistry key.Binding
	filter         key.Binding
	quit           key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Refresh"),
		),
		selectRegistry: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select"),
		),
		filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "Filter"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
	}
}

type Model struct {
	list    list.Model
	spinner spinner.Model

	registries []*Registry
	cursor     int

	width  int
	height int

	isFilterMode bool

	keys *listKeyMap
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

	return &Model{
		list:       l,
		spinner:    sp,
		registries: registries,
		keys:       newListKeyMap(),
	}
}

func (m *Model) Init() tea.Cmd {
	m.list.SetDelegate(&registryListDelegate{
		getSpinnerFunc: func() string {
			return m.spinner.View()
		},
	})

	cmds := []tea.Cmd{m.spinner.Tick, tea.WindowSize()}
	cmds = append(cmds, m.fetchRegistries()...)

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.selectRegistry) && m.registries[m.list.Index()].Status == registrystatus.Online:
			reg := m.registries[m.list.Index()]
			return nav.Navigate(registrydetails.NewRegistryDetails(&registrydetails.Registry{
				URL:            reg.URL,
				Username:       reg.Username,
				Password:       reg.Password,
				SupportsHTTPS3: reg.SupportsHTTP3,
			}, m))

		case key.Matches(msg, m.keys.refresh) && m.list.FilterState() != list.Filtering:
			return m, tea.Batch(m.fetchRegistries()...)

		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetHeight(m.height - 12)
		m.list.SetWidth(m.width - 5)

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

func (m *Model) View() string {
	s := &strings.Builder{}

	statusBar := ""

	out := lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.drawHeader(m.width-4, 5),
		ui.NewWindow(ui.WindowConfig{
			Width:       m.width - 4,
			Height:      m.height - 10,
			LeftTitle:   "Container Registries",
			RightTitle:  ui.CountKind(len(m.registries), "registry", "registries"),
			Content:     m.list.View(),
			BorderColor: lipgloss.Color("#FFFFFF"),
		}),
		statusBar,
	))

	fmt.Fprint(s, out)

	return s.String()
}

func (m *Model) drawHeader(width, height int) string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		ui.NewKeysWindow(ui.KeysWindowConfig{
			Width:  width - 50,
			Height: height,
			Keys: []key.Binding{
				key.NewBinding(key.WithHelp("↑/↓", "Navigate")),
				m.keys.selectRegistry,
				m.keys.filter,
				m.keys.refresh,
				m.keys.quit,
			},
		}),
		" ",
		drawStatusContent(height),
	)

	return out
}

func drawStatusContent(height int) string {
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

	return ui.NewWindow(ui.WindowConfig{
		Width:     18,
		Height:    height,
		LeftTitle: "Statuses",
		Content:   statusContainer.Render(statusContent),
	})
}

func (m *Model) fetchRegistries() []tea.Cmd {
	return lo.Map(m.registries, func(item *Registry, i int) tea.Cmd {
		item.Status = registrystatus.Loading
		return fetchRegistry(i, item)
	})
}
