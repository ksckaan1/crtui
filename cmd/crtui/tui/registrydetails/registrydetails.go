package registrydetails

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
	"github.com/ksckaan1/crtui/cmd/crtui/tui/tagdetails"
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

var _ tea.Model = (*Model)(nil)

type Model struct {
	// ui
	keysWindow         *ui.KeysWindow
	spinner            spinner.Model
	backModel          tea.Model
	repositoriesWindow *ui.Window
	repositoryListUI   list.Model
	tagsWindow         *ui.Window
	tagListUI          list.Model

	// Data
	registry  *Registry
	rc        *registryclient.RegistryClient
	startTime time.Time

	// Sizing
	width, height int

	// States
	activePaneIndex int
	// repositories
	isRepositoriesLoading bool
	repositoryList        []*Repository
	selectedRepository    *string
	// tags
	isTagsLoading bool
	tagList       []*Tag
	selectedTag   string

	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	status                 *ui.Status
}

func NewRegistryDetails(
	registry *Registry,
	backModel tea.Model,
	status *ui.Status,
) *Model {
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

	tagList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	tagList.SetShowTitle(false)
	tagList.SetShowHelp(false)
	tagList.DisableQuitKeybindings()
	tagList.SetShowStatusBar(false)

	return &Model{
		registry:           registry,
		rc:                 rc,
		keysWindow:         kw,
		backModel:          backModel,
		spinner:            sp,
		repositoriesWindow: ui.NewWindow(),
		repositoryListUI:   repositoryList,
		tagsWindow:         ui.NewWindow(),
		tagListUI:          tagList,

		isRepositoriesLoading: true,
		isTagsLoading:         false,

		selectedRepository:     new(""),
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),
		status:                 status,
	}
}

func (m *Model) Init() tea.Cmd {
	m.repositoryListUI.SetDelegate(&repositoryListDelegate{
		selectedRepository: m.selectedRepository,
	})

	m.tagListUI.SetDelegate(&tagListDelegate{})

	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.RequestWindowSize,
		m.fetchRepositoryList(),
	}

	if *m.selectedRepository != "" {
		cmds = append(cmds, m.fetchTagList())
	}

	m.startTime = time.Now()
	m.status.SetStatus(ui.Info, "Loading repository list...")

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Esc) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			m.tagListUI.FilterState() != list.Filtering:
			if m.activePaneIndex == 0 {
				m.status.SetStatus(ui.Empty, "")
				return nav.Navigate(m.backModel)
			} else {
				m.activePaneIndex = 0
			}

		case key.Matches(msg, keyMap.SelectItem) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			m.tagListUI.FilterState() != list.Filtering &&
			len(m.repositoryList) > 1:
			if m.activePaneIndex == 0 {
				*m.selectedRepository = m.repositoryList[m.repositoryListUI.GlobalIndex()].Name
				m.activePaneIndex = 1
				m.isTagsLoading = true
				m.status.SetStatus(ui.Info, "Loading tag list...")
				m.startTime = time.Now()

				return m, tea.Batch(m.fetchTagList(), tea.RequestWindowSize)
			}

			if m.activePaneIndex == 1 {
				m.selectedTag = m.tagList[m.tagListUI.GlobalIndex()].Name
				m.status.SetStatus(ui.Empty, "")

				return nav.Navigate(
					tagdetails.NewModel(
						m.rc,
						*m.selectedRepository,
						m.selectedTag,
						m,
						m.status,
					),
				)
			}

		case key.Matches(msg, keyMap.SwitchPane) &&
			*m.selectedRepository != "" &&
			m.repositoryListUI.FilterState() != list.Filtering:
			m.activePaneIndex = 1 - (m.activePaneIndex / 1)

		case key.Matches(msg, keyMap.Refresh) &&
			m.repositoryListUI.FilterState() != list.Filtering &&
			m.tagListUI.FilterState() != list.Filtering:

			if m.activePaneIndex == 0 && !m.isRepositoriesLoading {
				m.isRepositoriesLoading = true
				m.status.SetStatus(ui.Info, "Loading repository list...")
				m.startTime = time.Now()
				return m, tea.Batch(m.fetchRepositoryList(), tea.RequestWindowSize)
			}

			if m.activePaneIndex == 1 && !m.isTagsLoading {
				m.isTagsLoading = true
				m.status.SetStatus(ui.Info, "Loading tag list...")
				m.startTime = time.Now()
				return m, tea.Batch(m.fetchTagList(), tea.RequestWindowSize)
			}

		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.keysWindow.SetWidth(m.width - 29)
		m.repositoriesWindow.SetHeight(m.height - 8 - m.status.Height())
		if *m.selectedRepository == "" {
			m.repositoriesWindow.SetWidth(m.width - 2)
		} else {
			m.repositoriesWindow.SetWidth(m.width/2 - 1)
		}
		m.tagsWindow.SetHeight(m.height - 8 - m.status.Height())
		m.tagsWindow.SetWidth(m.width - lipgloss.Width(m.repositoriesWindow.View()) - 2)

	case repositoryListResult:
		m.isRepositoriesLoading = false
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

	case tagListResult:
		m.isTagsLoading = false
		m.tagListUI.ResetFilter()
		m.tagListUI.ResetSelected()
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

		if msg.err != nil {
			m.status.SetStatus(
				ui.Error,
				fmt.Sprintf("Failed to load tag list: %s", msg.err),
				time.Since(m.startTime),
			)
		} else {
			m.status.SetStatus(
				ui.OK,
				"Tag list loaded",
				time.Since(m.startTime),
			)
		}
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

func (m *Model) View() tea.View {
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
				m.drawContent(),
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
	)

	return out
}

func (m *Model) drawContent() string {
	content := []string{}

	content = append(content, m.drawRepositories())

	if *m.selectedRepository != "" {
		content = append(content,
			m.drawTags(),
		)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		content...,
	)
}

func (m *Model) drawRepositories() string {
	m.repositoriesWindow.SetBorderColor(nil)
	if m.activePaneIndex == 0 {
		m.repositoriesWindow.SetBorderColor(lipgloss.Color("#FFFFFF"))
	}

	repoLeftTitle := "Repositories"
	if m.repositoryListUI.IsFiltered() {
		repoLeftTitle += lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Render(" /" + m.repositoryListUI.FilterValue())
	}
	m.repositoriesWindow.SetLeftTitle(repoLeftTitle)

	m.repositoriesWindow.SetRightTitle(lo.Ternary(
		m.isRepositoriesLoading,
		m.spinner.View(),
		ui.CountKind(
			len(m.repositoryList),
			"repository",
			"repositories",
		),
	))

	m.repositoriesWindow.SetContent(func(width, height int) string {
		m.repositoryListUI.SetWidth(width)
		m.repositoryListUI.SetHeight(height)
		return m.repositoryListUI.View()
	})

	return m.repositoriesWindow.View()
}

func (m *Model) drawTags() string {
	m.tagsWindow.SetBorderColor(nil)
	if m.activePaneIndex == 1 {
		m.tagsWindow.SetBorderColor(lipgloss.Color("#FFFFFF"))
	}

	tagLeftTitle := "Tags"
	if m.tagListUI.IsFiltered() {
		tagLeftTitle += lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Render(" /" + m.tagListUI.FilterValue())
	}
	m.tagsWindow.SetLeftTitle(tagLeftTitle)

	m.tagsWindow.SetRightTitle(lo.Ternary(
		m.isTagsLoading,
		m.spinner.View(),
		ui.CountKind(
			len(m.tagList),
			"tag",
			"tags",
		),
	))

	m.tagsWindow.SetContent(func(width, height int) string {
		m.tagListUI.SetWidth(width)
		m.tagListUI.SetHeight(height)
		return m.tagListUI.View()
	})

	return m.tagsWindow.View()
}
