package tagdetails

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	rc         *registryclient.RegistryClient
	keys       *listKeyMap
	spinner    spinner.Model
	platformVP viewport.Model
	backModel  tea.Model

	width, height int

	host           string
	repositoryName string
	tagName        string
	tag            *models.Tag
	err            error

	isLoaded       bool
	activeTabIndex int
}

func NewModel(rc *registryclient.RegistryClient, repositoryName, tagName string, backModel tea.Model) *Model {
	sp := spinner.New()
	sp.Spinner = ui.LoadingSpinner

	platformVP := viewport.New(10, 10)

	host := strings.TrimPrefix(rc.BaseURL, "https://")
	host = strings.TrimPrefix(host, "http://")

	return &Model{
		rc:             rc,
		keys:           newListKeyMap(),
		spinner:        sp,
		platformVP:     platformVP,
		backModel:      backModel,
		host:           host,
		repositoryName: repositoryName,
		tagName:        tagName,
	}
}

func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.WindowSize(),
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

		case msg.String() == "left" && m.tag != nil:
			m.activeTabIndex--
			if m.activeTabIndex < 0 {
				m.activeTabIndex = len(m.tag.Platforms) - 1
			}
			m.drawPlatformContent()

		case msg.String() == "right" && m.tag != nil:
			m.activeTabIndex++
			if m.activeTabIndex >= len(m.tag.Platforms) {
				m.activeTabIndex = 0
			}
			m.drawPlatformContent()

		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.platformVP.Width = m.width - 4
		m.platformVP.Height = m.height - 17
		m.drawPlatformContent()

	case tagMsg:
		m.tag = msg.tag
		m.err = msg.err
		m.isLoaded = true
		m.drawPlatformContent()
	}

	cmds := []tea.Cmd{}

	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	var vpCmd tea.Cmd
	m.platformVP, vpCmd = m.platformVP.Update(msg)

	cmds = append(cmds, spinnerCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := &strings.Builder{}

	statusBar := ""

	rightTitle := ""

	if !m.isLoaded {
		rightTitle = m.spinner.View()
	}

	if m.tag == nil {
		return lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.drawHeader(m.width-4, 5),
			ui.NewWindow(ui.WindowConfig{
				Width:      m.width - 4,
				Height:     m.height - 10,
				RightTitle: rightTitle,
				Content:    m.drawTop(m.width-4, 5),
			}),
			statusBar,
		))
	}

	out := lipgloss.NewStyle().Padding(1).Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.drawHeader(m.width-4, 5),
		ui.NewWindow(ui.WindowConfig{
			Width:      m.width - 4,
			Height:     3,
			RightTitle: rightTitle,
			Content:    m.drawTop(m.width-4, 5),
		}),
		ui.NewTabbedWindow(ui.TabbedWindowConfig{
			Width:          m.width - 4,
			Height:         m.height - 16,
			Content:        m.platformVP.View(),
			ActiveTabIndex: m.activeTabIndex,
			Tabs: lo.Map(m.tag.Platforms, func(item models.Platform, _ int) string {
				return fmt.Sprintf("%s/%s", item.Config.OS, item.Config.Architecture)
			}),
		}),
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

func (m *Model) drawTop(width, height int) string {
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

func (m *Model) drawPlatformContent() {
	if !m.isLoaded {
		return
	}

	strs := []string{}

	labels := m.drawLabels()
	if labels != "" {
		strs = append(strs, labels)
	}

	envs := m.drawEnvs()
	if envs != "" {
		strs = append(strs, envs)
	}

	cmds := m.drawCmd()
	if cmds != "" {
		strs = append(strs, cmds)
	}

	entrypoint := m.drawEntrypoint()
	if entrypoint != "" {
		strs = append(strs, entrypoint)
	}

	workingDir := m.drawWorkingDir()
	if workingDir != "" {
		strs = append(strs, workingDir)
	}

	user := m.drawUser()
	if user != "" {
		strs = append(strs, user)
	}

	strs = append(strs,
		m.drawLayers(),
		m.drawHistory(),
	)

	rootFS := m.drawRootFS()
	if rootFS != "" {
		strs = append(strs, rootFS)
	}

	m.platformVP.SetContent(lipgloss.JoinVertical(lipgloss.Top, strs...))
}

func (m *Model) drawGrid(fields [][2]any, minElementWidth, rowWidth int) string {
	n := len(fields)
	if n == 0 {
		return ""
	}

	maxCn := max(rowWidth/minElementWidth, 1)

	cn := min(n, maxCn) // column count

	cw := rowWidth / cn // column width

	cellStyle := lipgloss.NewStyle().
		Width(cw).
		Align(lipgloss.Left).
		MarginRight(1).
		MarginBottom(1)

	var rows []string
	var currentRow []string

	for i, field := range fields {
		renderedCell := cellStyle.Render(m.drawField(field[0], field[1], cw))
		currentRow = append(currentRow, renderedCell)

		if (i+1)%cn == 0 || i+1 == n {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *Model) drawField(title, content any, width int) string {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false).
		BorderLeft(true).
		PaddingLeft(1).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().
				Faint(true).
				Width(width).
				Render(fmt.Sprintf("%v", title)),
			lipgloss.NewStyle().
				Width(width).
				Render(fmt.Sprintf("%v", content)),
		))
}

func humanReadableSize(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}
