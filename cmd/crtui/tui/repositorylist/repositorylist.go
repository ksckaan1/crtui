package repositorylist

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/figlet"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/taglist"
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

var _ tea.Model = (*RepositoryListScreenModel)(nil)

type RepositoryListScreenModel struct {
	// ui
	keysWindow         *ui.KeysWindow
	spinner            spinner.Model
	backModel          tea.Model
	repositoriesWindow *ui.Window
	repositoryListUI   list.Model

	// Data
	registry  *Registry
	rc        *registryclient.RegistryClient
	startTime time.Time

	// Sizing
	width, height int

	isLoading          bool
	repositoryList     []*Repository
	selectedRepository *string

	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	status                 *ui.Status
}

func NewRepositoryListScreenModel(
	registry *Registry,
	backModel tea.Model,
	status *ui.Status,
) *RepositoryListScreenModel {
	rc := registryclient.New(
		registry.URL,
		registry.Username,
		registry.Password,
	)
	rc.SetTransport(registry.SupportsHTTPS3)

	sp := spinner.New()
	sp.Spinner = ui.LoadingSpinner

	kw := ui.NewKeysWindow()
	kw.SetHeight(5)
	kw.SetKeyMap(keyMap)

	repositoryList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	repositoryList.SetShowTitle(false)
	repositoryList.SetShowHelp(false)
	repositoryList.DisableQuitKeybindings()
	repositoryList.SetShowStatusBar(false)

	return &RepositoryListScreenModel{
		registry:           registry,
		rc:                 rc,
		keysWindow:         kw,
		backModel:          backModel,
		spinner:            sp,
		repositoriesWindow: ui.NewWindow(),
		repositoryListUI:   repositoryList,

		isLoading: true,

		selectedRepository:     new(""),
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),
		status:                 status,
	}
}

func (m *RepositoryListScreenModel) Init() tea.Cmd {
	m.repositoryListUI.SetDelegate(&repositoryListDelegate{
		selectedRepository: m.selectedRepository,
	})

	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.RequestWindowSize,
		m.fetchRepositoryList(),
	}

	m.startTime = time.Now()
	m.status.SetStatus(ui.Info, "Loading repository list...")

	return tea.Batch(cmds...)
}

func (m *RepositoryListScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Esc) &&
			m.repositoryListUI.FilterState() != list.Filtering:
			m.status.SetStatus(ui.Empty, "")
			return nav.Navigate(m.backModel)

		case key.Matches(msg, keyMap.Select) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			len(m.repositoryList) > 0:
			return nav.Navigate(
				taglist.NewTagListScreenModel(
					m.rc,
					m,
					m.status,
					m.repositoryList[m.repositoryListUI.GlobalIndex()].Name,
				),
			)

		case key.Matches(msg, keyMap.Delete) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			len(m.repositoryList) > 0:

			return nav.Navigate(
				NewDeleteTagsPopup(
					m.rc,
					m,
					m.status,
					m.repositoryList[m.repositoryListUI.GlobalIndex()].Name,
					nil,
					m.spinner,
				),
			)

		case key.Matches(msg, keyMap.Refresh) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			!m.isLoading:

			m.isLoading = true
			m.status.SetStatus(ui.Info, "Loading repository list...")
			m.startTime = time.Now()
			return m, tea.Batch(m.fetchRepositoryList(), tea.RequestWindowSize)

		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.UpdateSize(msg)

	case repositoryListResult:
		m.isLoading = false
		m.repositoryListUI.ResetFilter()
		m.repositoryListUI.ResetSelected()

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

		if msg.err != nil {
			m.status.SetStatus(
				ui.Error,
				fmt.Sprintf("Failed to load repository list: %s", msg.err),
				time.Since(m.startTime),
			)
		} else {
			m.status.SetStatus(
				ui.OK,
				"Repository list loaded",
				time.Since(m.startTime),
			)
		}
	}

	cmds := []tea.Cmd{}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	var repositoryListCmd tea.Cmd
	m.repositoryListUI, repositoryListCmd = m.repositoryListUI.Update(msg)
	cmds = append(cmds, repositoryListCmd)

	return m, tea.Batch(cmds...)
}

func (m *RepositoryListScreenModel) UpdateSize(msg tea.WindowSizeMsg) {
	m.width, m.height = msg.Width, msg.Height
	m.keysWindow.SetWidth(m.width - 29)
	m.repositoriesWindow.SetWidth(m.width - 2)
	m.repositoriesWindow.SetHeight(m.height - 8 - m.status.Height())
}

func (m *RepositoryListScreenModel) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	sb := &strings.Builder{}

	out := lipgloss.NewStyle().
		Padding(1).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.drawHeader(),
				m.drawRepositories(),
				m.status.View(),
			),
		)

	fmt.Fprint(sb, out)

	v.SetContent(sb.String())

	return v
}

func (m *RepositoryListScreenModel) drawHeader() string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		m.keysWindow.View(),
	)

	return out
}

func (m *RepositoryListScreenModel) drawRepositories() string {
	repoLeftTitle := "Repositories"
	if m.repositoryListUI.IsFiltered() {
		repoLeftTitle += lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Render(" /" + m.repositoryListUI.FilterValue())
	}
	m.repositoriesWindow.SetLeftTitle(repoLeftTitle)

	m.repositoriesWindow.SetRightTitle(lo.Ternary(
		m.isLoading,
		m.spinner.View(),
		ui.CountKind(
			len(m.repositoryList),
			"repository",
			"repositories",
		),
	))

	m.repositoriesWindow.SetContent(func(width, height int) string {
		if len(m.repositoryList) == 0 && m.isLoading {
			return "Loading..."
		}

		if len(m.repositoryList) == 0 {
			return "No repositories found"
		}

		m.repositoryListUI.SetWidth(width)
		m.repositoryListUI.SetHeight(height)
		return m.repositoryListUI.View()
	})

	return m.repositoriesWindow.View()
}
