package registrydetails

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/figlet"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
	"github.com/samber/lo"
)

type Registry struct {
	URL            string
	Username       string
	Password       string
	SupportsHTTPS3 bool
}

type listKeyMap struct {
	refresh    key.Binding
	selectItem key.Binding
	filter     key.Binding
	esc        key.Binding
	quit       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Refresh"),
		),
		selectItem: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select item"),
		),
		filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "Filter"),
		),
		esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Back"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
	}
}

var _ tea.Model = (*Model)(nil)

type Model struct {
	// ui
	spinner          spinner.Model
	keys             *listKeyMap
	backModel        tea.Model
	repositoryListUI list.Model
	tagListUI        list.Model

	// Data
	registry *Registry
	rc       *registryclient.RegistryClient

	// Sizing
	width, height int

	// States
	activePaneIndex int
	// repositories
	isRepositoriesLoading bool
	repositoryList        []*Repository
	selectedRepository    string
	// tags
	isTagsLoading bool
	tagList       []*Tag
}

func NewRegistryDetails(registry *Registry, backModel tea.Model) *Model {
	rc := registryclient.New(
		registry.URL,
		registry.Username,
		registry.Password,
	)

	rc.SetTransport(registry.SupportsHTTPS3)

	sp := spinner.New()
	sp.Spinner = spinner.Spinner{
		// white \033[37m
		// magenta \033[35m
		// reset \033[0m
		Frames: []string{
			"\033[35ml\033[37moading\033[0m",         // [l]oading
			"\033[37ml\033[35mo\033[37mading\033[0m", // l[o]ading
			"\033[37mlo\033[35ma\033[37mding\033[0m", // lo[a]ding
			"\033[37mloa\033[35md\033[37ming\033[0m", // loa[d]ing
			"\033[37mload\033[35mi\033[37mng\033[0m", // load[i]ng
			"\033[37mloadi\033[35mn\033[37mg\033[0m", // loadi[n]g
			"\033[37mloadin\033[35mg\033[0m",         // loadin[g]
		},
		FPS: time.Second / 7,
	}

	repositoryList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	repositoryList.SetShowTitle(false)
	repositoryList.SetShowHelp(false)
	repositoryList.DisableQuitKeybindings()
	repositoryList.SetShowStatusBar(false)

	tagList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	tagList.SetShowTitle(false)
	tagList.SetShowHelp(false)
	tagList.DisableQuitKeybindings()
	tagList.SetShowStatusBar(false)

	return &Model{
		registry:         registry,
		rc:               rc,
		backModel:        backModel,
		keys:             newListKeyMap(),
		spinner:          sp,
		repositoryListUI: repositoryList,
		tagListUI:        tagList,

		isRepositoriesLoading: true,
		isTagsLoading:         false,
	}
}

func (m *Model) Init() tea.Cmd {
	m.repositoryListUI.SetDelegate(&repositoryListDelegate{})
	m.tagListUI.SetDelegate(&tagListDelegate{})

	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.WindowSize(),
		m.fetchRepositoryList(),
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.esc) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			m.tagListUI.FilterState() != list.Filtering:
			if m.activePaneIndex == 0 {
				return nav.Navigate(m.backModel)
			} else {
				m.activePaneIndex = 0
			}

		case key.Matches(msg, m.keys.selectItem) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			len(m.repositoryList) > 1:
			if m.activePaneIndex == 0 {
				m.selectedRepository = m.repositoryList[m.repositoryListUI.Index()].Name
				m.activePaneIndex = 1
				m.isTagsLoading = true

				return m, m.fetchTagList()
			}

		case msg.String() == "left" && m.selectedRepository != "":
			if m.activePaneIndex != 0 {
				m.activePaneIndex = 0
			}

		case msg.String() == "right" && m.selectedRepository != "":
			if m.activePaneIndex != 1 {
				m.activePaneIndex = 1
			}

		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.repositoryListUI.SetHeight(m.height - 12)
		m.repositoryListUI.SetWidth(m.width - 5)
		m.tagListUI.SetHeight(m.height - 12)
		m.tagListUI.SetWidth(m.width - 5)

	case repositoryListResult:
		m.isRepositoriesLoading = false
		m.repositoryList = lo.Map(msg.repositoryList, func(item string, _ int) *Repository {
			return &Repository{
				Name: item,
			}
		})
		m.repositoryListUI.SetItems(lo.Map(m.repositoryList, func(item *Repository, _ int) list.Item {
			return &Repository{
				Name: item.Name,
			}
		}))

	case tagListResult:
		m.isTagsLoading = false
		m.tagList = lo.Map(msg.tagList, func(item string, _ int) *Tag {
			return &Tag{
				Name: item,
			}
		})
		m.tagListUI.SetItems(lo.Map(m.tagList, func(item *Tag, _ int) list.Item {
			return &Tag{
				Name: item.Name,
			}
		}))
	}

	cmds := []tea.Cmd{}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	if m.activePaneIndex == 0 {
		var repositoryListCmd tea.Cmd
		m.repositoryListUI, repositoryListCmd = m.repositoryListUI.Update(msg)
		cmds = append(cmds, repositoryListCmd)
	}

	if m.activePaneIndex == 1 {
		var tagListCmd tea.Cmd
		m.tagListUI, tagListCmd = m.tagListUI.Update(msg)
		cmds = append(cmds, tagListCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := &strings.Builder{}

	statusBar := ""

	out := lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.drawHeader(m.width-4, 5),
		m.drawContent(m.width-4, m.height-10),
		statusBar,
	))

	fmt.Fprint(sb, out)

	return sb.String()
}

func (m *Model) drawHeader(width, height int) string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		ui.NewKeysWindow(ui.KeysWindowConfig{
			Width:  width - 29,
			Height: height,
			Keys: []key.Binding{
				key.NewBinding(key.WithHelp("↑ ↓", "Navigate")),
				key.NewBinding(key.WithHelp("← →", "Switch Pane")),
				m.keys.selectItem,
				m.keys.refresh,
				m.keys.filter,
				m.keys.esc,
				m.keys.quit,
			},
		}),
	)

	return out
}

func (m *Model) drawContent(width, height int) string {
	paneWidth := (width / 2) - 1
	if m.selectedRepository == "" {
		paneWidth = width
	}

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.NewWindow(ui.WindowConfig{
			Width:     paneWidth,
			Height:    height,
			LeftTitle: "Repositories",
			RightTitle: lo.Ternary(
				m.isRepositoriesLoading,
				m.spinner.View(),
				ui.CountKind(
					len(m.repositoryList),
					"repository",
					"repositories",
				),
			),
			Content: m.drawRepositoryList(),
			BorderColor: lo.Ternary(
				m.activePaneIndex == 0,
				lipgloss.Color("#FFFFFF"),
				"",
			),
		}),
		lo.Ternary(
			m.selectedRepository != "",
			ui.NewWindow(ui.WindowConfig{
				Width:     paneWidth,
				Height:    height,
				LeftTitle: "Tags",
				RightTitle: lo.Ternary(
					m.isTagsLoading,
					m.spinner.View(),
					ui.CountKind(len(m.tagList), "tag", "tags"),
				),
				Content: m.drawTagList(),
				BorderColor: lo.Ternary(
					m.activePaneIndex == 1,
					lipgloss.Color("#FFFFFF"),
					"",
				),
			}),
			"",
		),
	)

	return out
}

func (m *Model) drawRepositoryList() string {
	if !m.isRepositoriesLoading && len(m.repositoryList) == 0 {
		return "No repository found!"
	}

	if m.isRepositoriesLoading && len(m.repositoryList) == 0 {
		return ""
	}

	return m.repositoryListUI.View()
}

func (m *Model) drawTagList() string {
	if !m.isTagsLoading && len(m.tagList) == 0 {
		return "No tag found!"
	}

	if m.isTagsLoading && len(m.tagList) == 0 {
		return ""
	}

	return m.tagListUI.View()
}
