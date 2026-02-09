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
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/samber/lo"
)

var _ tea.Model = (*Model)(nil)

type listKeyMap struct {
	quit    key.Binding
	refresh key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Refresh"),
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

	cmds := []tea.Cmd{m.spinner.Tick}
	cmds = append(cmds, m.fetchRegistries()...)

	return tea.Batch(cmds...)
}

func (m *Model) fetchRegistries() []tea.Cmd {
	return lo.Map(m.registries, func(item *Registry, i int) tea.Cmd {
		item.Status = "loading"
		return fetchRegistry(i, item)
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.refresh) && m.list.FilterState() != list.Filtering:
			return m, tea.Batch(m.fetchRegistries()...)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetHeight(m.height - 12)
		m.list.SetWidth(m.width - 5)

	case registryResult:
		m.registries[msg.index].Status = msg.status
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
		drawHeader(m.width-4, 5),
		ui.NewWindow(ui.WindowConfig{
			Width:      m.width - 4,
			Height:     m.height - 10,
			LeftTitle:  "Container Registries",
			RightTitle: fmt.Sprintf("%d registries", len(m.registries)),
			Content:    m.list.View(),
		}),
		statusBar,
	))

	fmt.Fprint(s, out)

	return s.String()
}

func drawHeader(width, height int) string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Foreground(primaryColor).Render(figlet.Figlet),
		" ",
		ui.NewWindow(ui.WindowConfig{
			Width:     width - 50,
			Height:    height,
			LeftTitle: "Keys",
			Content:   drawHelpContent(),
		}),
		" ",
		ui.NewWindow(ui.WindowConfig{
			Width:     18,
			Height:    height,
			LeftTitle: "Statuses",
			Content:   drawStatusContent(),
		}),
	)

	return out
}

func drawHelpContent() string {
	helpContainer := lipgloss.NewStyle().
		PaddingLeft(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Width(6)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	arrowKeys := fmt.Sprintf("%s %s", keyStyle.Render("↑/↓"), descStyle.Render("Navigate"))
	enterKey := fmt.Sprintf("%s %s", keyStyle.Render("enter"), descStyle.Render("Select"))
	quitKey := fmt.Sprintf("%s %s", keyStyle.Render("ctrl-c"), descStyle.Render("Quit"))
	filterKey := fmt.Sprintf("%s %s", keyStyle.Render("/"), descStyle.Render("Filter"))
	refreshKey := fmt.Sprintf("%s %s", keyStyle.Render("r"), descStyle.Render("Refresh"))

	helpContent := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			arrowKeys,
			enterKey,
			filterKey,
			refreshKey,
		),
		"  ",
		lipgloss.JoinVertical(
			lipgloss.Left,
			quitKey,
		),
	)

	return helpContainer.Render(helpContent)
}

func drawStatusContent() string {
	statusContainer := lipgloss.NewStyle().
		PaddingLeft(2)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))

	onlineStatus := fmt.Sprintf("%s %s", onlineDot, descStyle.Render("Online"))
	offlineStatus := fmt.Sprintf("%s %s", offlineDot, descStyle.Render("Offline"))
	unauthStatus := fmt.Sprintf("%s %s", unauthDot, descStyle.Render("Unauthorized"))

	statusContent := lipgloss.JoinVertical(
		lipgloss.Left,
		onlineStatus,
		offlineStatus,
		unauthStatus,
	)

	return statusContainer.Render(statusContent)
}
