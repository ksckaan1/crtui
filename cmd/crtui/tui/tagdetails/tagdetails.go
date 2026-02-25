package tagdetails

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/figlet"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/nav"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
	"github.com/samber/lo"
)

type listKeyMap struct {
	refresh key.Binding
	esc     key.Binding
	quit    key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "Refresh"),
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
	rc              *registryclient.RegistryClient
	keys            *listKeyMap
	spinner         spinner.Model
	backModel       tea.Model
	platformsWindow *ui.TabbedWindow

	width, height int

	host           string
	repositoryName string
	tagName        string
	tag            *models.Tag
	err            error

	isLoaded bool
}

func NewModel(rc *registryclient.RegistryClient, repositoryName, tagName string, backModel tea.Model) *Model {
	sp := spinner.New()
	sp.Spinner = ui.LoadingSpinner

	pw := ui.NewTabbedWindow()

	host := strings.TrimPrefix(rc.BaseURL, "https://")
	host = strings.TrimPrefix(host, "http://")

	return &Model{
		rc:        rc,
		keys:      newListKeyMap(),
		spinner:   sp,
		backModel: backModel,

		platformsWindow: pw,
		host:            host,
		repositoryName:  repositoryName,
		tagName:         tagName,
	}
}

func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.RequestWindowSize,
		m.fetchTag(),
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.esc):
			return nav.Navigate(m.backModel)

		case key.Matches(msg, m.keys.refresh):
			m.isLoaded = false
			return m, m.fetchTag()

		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.platformsWindow.SetWidth(m.width - 2)
		m.platformsWindow.SetHeight(m.height - 14)

	case tagMsg:
		m.tag = msg.tag
		m.err = msg.err
		m.isLoaded = true
		if m.tag != nil {
			m.platformsWindow.SetTabs(lo.Map(m.tag.Platforms, func(item models.Platform, _ int) string {
				return fmt.Sprintf("%s/%s", item.Config.OS, item.Config.Architecture)
			}))
		}
	}

	cmds := []tea.Cmd{}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	var platformsCmd tea.Cmd
	m.platformsWindow, platformsCmd = m.platformsWindow.Update(msg)

	cmds = append(cmds, spinnerCmd, platformsCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	sb := &strings.Builder{}

	statusBar := "STATUS"

	rightTitle := ""

	if !m.isLoaded {
		rightTitle = m.spinner.View()
	}

	if m.tag == nil {
		v.SetContent(lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.drawHeader(m.width-2, 5),
			ui.NewWindow(ui.WindowConfig{
				Width:      m.width - 2,
				Height:     m.height - 8,
				RightTitle: rightTitle,
				Content:    m.drawTop(),
			}),
			statusBar,
		)))
		return v
	}

	m.platformsWindow.SetContent(m.drawPlatforms)

	out := lipgloss.NewStyle().
		Padding(1).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.drawHeader(m.width-2, 5),
			ui.NewWindow(ui.WindowConfig{
				Width:      m.width - 2,
				Height:     3,
				RightTitle: rightTitle,
				Content:    m.drawTop(),
			}),
			m.platformsWindow.View(),
			statusBar,
		))

	fmt.Fprint(sb, out)

	v.SetContent(sb.String())

	return v
}

func (m *Model) drawHeader(width, height int) string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		ui.NewKeysWindow(ui.KeysWindowConfig{
			Width:  width - 27,
			Height: height,
			Keys: []key.Binding{
				key.NewBinding(key.WithHelp("↑ ↓", "Scroll")),
				key.NewBinding(key.WithHelp("← →", "Switch Tabs")),
				m.keys.refresh,
				m.keys.esc,
				m.keys.quit,
			},
		}),
	)

	return out
}

func (m *Model) drawTop() string {
	if m.tag == nil {
		return ""
	}

	totalSizeBytes := lo.SumBy(m.tag.Platforms, func(item models.Platform) int {
		return lo.SumBy(item.Layers, func(layer models.Layer) int {
			return layer.Size
		})
	})

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().
			Width(m.width/2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Render("Name:"),
					lipgloss.NewStyle().
						Bold(true).
						Render(fmt.Sprintf(
							"%s:%s",
							m.repositoryName,
							m.tagName,
						)),
				),
			),
		lipgloss.NewStyle().
			Width(m.width/2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lipgloss.NewStyle().
						Faint(true).
						Render("Total Size:"),
					lipgloss.NewStyle().
						Bold(true).
						Render(fmt.Sprintf(
							"%s (%dB)",
							humanReadableSize(totalSizeBytes),
							totalSizeBytes,
						)),
				),
			),
	)
}

func (m *Model) drawPlatforms(width, height, activeTabIndex int) string {
	activePlatform := m.tag.Platforms[activeTabIndex]

	strs := []string{}

	labels := m.drawLabels(activePlatform, width)
	if labels != "" {
		strs = append(strs, labels)
	}

	envs := m.drawEnvs(activePlatform, width)
	if envs != "" {
		strs = append(strs, envs)
	}

	cmds := m.drawCmd(activePlatform, width)
	if cmds != "" {
		strs = append(strs, cmds)
	}

	entrypoint := m.drawEntrypoint(activePlatform, width)
	if entrypoint != "" {
		strs = append(strs, entrypoint)
	}

	workingDir := m.drawWorkingDir(activePlatform, width)
	if workingDir != "" {
		strs = append(strs, workingDir)
	}

	user := m.drawUser(activePlatform, width)
	if user != "" {
		strs = append(strs, user)
	}

	strs = append(strs,
		m.drawLayers(activePlatform, width),
		m.drawHistory(activePlatform, width),
	)

	rootFS := m.drawRootFS(activePlatform, width)
	if rootFS != "" {
		strs = append(strs, rootFS)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		strs...,
	)
}
