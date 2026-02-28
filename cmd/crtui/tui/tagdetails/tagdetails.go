package tagdetails

import (
	"fmt"
	"strings"
	"time"

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

var _ tea.Model = (*TagDetailsScreenModel)(nil)

type TagDetailsScreenModel struct {
	rc              *registryclient.RegistryClient
	spinner         spinner.Model
	back            tea.Model
	keysWindow      *ui.KeysWindow
	heroWindow      *ui.Window
	platformsWindow *ui.TabbedWindow
	status          *ui.Status

	width, height int

	pullCommand    string
	repositoryName string
	tagName        string
	tag            *models.Tag
	err            error

	isLoaded bool

	timeStart time.Time

	minTerminalSizeWarning *ui.MinTerminalSizeWarning
}

func NewTagDetailsScreenModel(
	rc *registryclient.RegistryClient,
	repositoryName,
	tagName string,
	back tea.Model,
	status *ui.Status,
) *TagDetailsScreenModel {
	sp := spinner.New()
	sp.Spinner = ui.LoadingSpinner

	hw := ui.NewWindow()
	pw := ui.NewTabbedWindow()

	host := strings.TrimPrefix(rc.BaseURL, "https://")
	host = strings.TrimPrefix(host, "http://")

	kw := ui.NewKeysWindow()
	kw.SetHeight(5)
	kw.SetKeyMap(keyMap)

	return &TagDetailsScreenModel{
		rc:         rc,
		spinner:    sp,
		back:       back,
		keysWindow: kw,
		status:     status,

		heroWindow:             hw,
		platformsWindow:        pw,
		pullCommand:            fmt.Sprintf("docker pull %s/%s:%s", host, repositoryName, tagName),
		repositoryName:         repositoryName,
		tagName:                tagName,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(82, 24),

		timeStart: time.Now(),
	}
}

func (m *TagDetailsScreenModel) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.RequestWindowSize,
		m.fetchTag(),
	}

	m.status.SetStatus(ui.Info, "Loading tag details...")

	return tea.Batch(cmds...)
}

func (m *TagDetailsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.DeleteTag):
			return nav.Navigate(
				NewDeleteTagPopupModel(
					m.rc,
					m,
					m.status,
					m.repositoryName,
					m.tagName,
				),
			)

		case key.Matches(msg, keyMap.Esc):
			m.status.SetStatus(ui.Empty, "")
			return nav.Navigate(m.back)

		case key.Matches(msg, keyMap.Refresh):
			m.isLoaded = false
			m.status.SetStatus(ui.Info, "Refreshing tag details...")
			m.timeStart = time.Now()
			return m, m.fetchTag()

		case key.Matches(msg, keyMap.CopyPullCommand):
			m.status.SetStatus(
				ui.Info,
				fmt.Sprintf("Pull command copied (%s)", m.pullCommand),
			)
			return m, tea.SetClipboard(m.pullCommand)

		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.UpdateSize(msg)

	case tagMsg:
		m.tag = msg.tag
		m.err = msg.err
		m.isLoaded = true

		took := time.Since(m.timeStart)

		if m.tag != nil {
			m.platformsWindow.SetTabs(lo.Map(m.tag.Platforms, func(item models.Platform, _ int) string {
				return fmt.Sprintf("%s/%s", item.Config.OS, item.Config.Architecture)
			}))
			m.status.SetStatus(ui.OK, "Tag details loaded", took)
		}

		if m.err != nil {
			m.status.SetStatus(ui.Error, fmt.Sprintf("Failed to load tag details: %s", msg.err.Error()), took)
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

func (m *TagDetailsScreenModel) UpdateSize(msg tea.WindowSizeMsg) {
	m.width, m.height = msg.Width, msg.Height
	m.keysWindow.SetWidth(m.width - 29)
	m.heroWindow.SetWidth(m.width - 2)
	m.platformsWindow.SetWidth(m.width - 2)
	m.platformsWindow.SetHeight(m.height - 13 - lipgloss.Height(m.status.View()))
}

func (m *TagDetailsScreenModel) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	wrn := m.minTerminalSizeWarning.View()
	if wrn != "" {
		v.SetContent(wrn)
		return v
	}

	sb := &strings.Builder{}

	m.heroWindow.SetRightTitle(lo.Ternary(!m.isLoaded, m.spinner.View(), ""))

	if m.tag == nil {
		m.heroWindow.SetHeight(m.height - 8 - m.status.Height())

		v.SetContent(lipgloss.NewStyle().
			Padding(1).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					m.drawHeader(),
					m.heroWindow.View(),
					m.status.View(),
				),
			),
		)

		return v
	}

	m.heroWindow.SetHeight(3)
	m.heroWindow.SetContent(m.drawTop)

	m.platformsWindow.SetContent(m.drawPlatforms)

	out := lipgloss.NewStyle().
		Padding(1).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.drawHeader(),
			m.heroWindow.View(),
			m.platformsWindow.View(),
			m.status.View(),
		))

	fmt.Fprint(sb, out)

	v.SetContent(sb.String())

	return v
}

func (m *TagDetailsScreenModel) drawHeader() string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		m.keysWindow.View(),
	)

	return out
}

func (m *TagDetailsScreenModel) drawTop(width, _ int) string {
	totalSizeBytes := lo.SumBy(m.tag.Platforms, func(item models.Platform) int {
		return lo.SumBy(item.Layers, func(layer models.Layer) int {
			return layer.Size
		})
	})

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().
			Width(width/2).
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
			Width(width/2).
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

func (m *TagDetailsScreenModel) drawPlatforms(width, height, activeTabIndex int) string {
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
