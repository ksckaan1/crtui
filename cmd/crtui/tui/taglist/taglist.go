package taglist

import (
	"fmt"
	"slices"
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

var _ tea.Model = (*TagListScreenModel)(nil)

type TagListScreenModel struct {
	// ui
	keysWindow             *ui.KeysWindow
	spinner                spinner.Model
	backModel              tea.Model
	tagsWindow             *ui.Window
	tagListUI              list.Model
	selectedRepositoryName string

	// Data
	rc        *registryclient.RegistryClient
	startTime time.Time

	// Sizing
	width, height int

	// States
	isLoading  bool
	tagList    []*Tag
	markedTags *[]*Tag

	minTerminalSizeWarning *ui.MinTerminalSizeWarning
	status                 *ui.Status
}

func NewTagListScreenModel(
	rc *registryclient.RegistryClient,
	backModel tea.Model,
	status *ui.Status,
	repositoryName string,
) *TagListScreenModel {
	sp := spinner.New()
	sp.Spinner = ui.LoadingSpinner

	kw := ui.NewKeysWindow()
	kw.SetHeight(5)
	kw.SetKeyMap(keyMap)

	tagList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	tagList.SetShowTitle(false)
	tagList.SetShowHelp(false)
	tagList.DisableQuitKeybindings()
	tagList.SetShowStatusBar(false)

	return &TagListScreenModel{
		rc:                     rc,
		keysWindow:             kw,
		backModel:              backModel,
		spinner:                sp,
		tagsWindow:             ui.NewWindow(),
		tagListUI:              tagList,
		isLoading:              true,
		selectedRepositoryName: repositoryName,
		minTerminalSizeWarning: ui.NewMinTerminalSizeWarning(60, 24),
		status:                 status,
		markedTags:             new([]*Tag),
	}
}

func (m *TagListScreenModel) Init() tea.Cmd {
	m.tagListUI.SetDelegate(&tagListDelegate{
		markedTags: m.markedTags,
	})

	cmds := []tea.Cmd{
		m.spinner.Tick,
		tea.RequestWindowSize,
		m.fetchTagList(),
	}

	m.startTime = time.Now()
	m.status.SetStatus(ui.Info, "Loading tag list...")

	return tea.Batch(cmds...)
}

func (m *TagListScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Esc) &&
			m.tagListUI.FilterState() != list.Filtering:
			if len(*m.markedTags) > 0 {
				*m.markedTags = []*Tag{}
				m.status.SetStatus(ui.Info, "Marked tags cleared")
				return m, nil
			}

			if m.tagListUI.IsFiltered() {
				m.tagListUI.ResetFilter()
				m.status.SetStatus(ui.Info, "Filter reset")
				return m, nil
			}

			m.status.SetStatus(ui.Empty, "")

			return nav.Navigate(m.backModel)

		// SELECT TAG
		case key.Matches(msg, keyMap.Select) &&
			m.tagListUI.FilterState() != list.Filtering &&
			len(m.tagList) > 0:
			m.status.SetStatus(ui.Empty, "")

			return nav.Navigate(
				tagdetails.NewTagDetailsScreenModel(
					m.rc,
					m.selectedRepositoryName,
					m.tagList[m.tagListUI.GlobalIndex()].Name,
					m,
					m.status,
				),
			)

		// MARK TAGS
		case key.Matches(msg, keyMap.Mark) &&
			m.tagListUI.FilterState() != list.Filtering:
			m.toggleMarkedTag()

		// DELETE TAGS
		case key.Matches(msg, keyMap.Delete) &&
			m.tagListUI.FilterState() != list.Filtering &&
			len(m.tagList) != 0:
			tagNames := lo.Map(*m.markedTags, func(item *Tag, _ int) string {
				return item.Name
			})

			if len(tagNames) == 0 {
				tagNames = append(tagNames, m.tagList[m.tagListUI.GlobalIndex()].Name)
			}

			m.status.SetStatus(ui.Empty, "")

			*m.markedTags = []*Tag{}

			return nav.Navigate(NewDeleteTagsPopup(
				m.rc,
				m,
				m.status,
				m.selectedRepositoryName,
				tagNames,
				m.spinner,
			))

		// REFRESH
		case key.Matches(msg, keyMap.Refresh) &&
			m.tagListUI.FilterState() != list.Filtering:

			m.isLoading = true
			m.status.SetStatus(ui.Info, "Loading tag list...")
			m.startTime = time.Now()
			return m, tea.Batch(m.fetchTagList(), tea.RequestWindowSize)

		// QUIT
		case key.Matches(msg, keyMap.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.UpdateSize(msg)

	case tagListResult:
		m.isLoading = false
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

	var tagListCmd tea.Cmd
	m.tagListUI, tagListCmd = m.tagListUI.Update(msg)
	cmds = append(cmds, tagListCmd)

	return m, tea.Batch(cmds...)
}

func (m *TagListScreenModel) UpdateSize(msg tea.WindowSizeMsg) {
	m.width, m.height = msg.Width, msg.Height
	m.keysWindow.SetWidth(m.width - 29)
	m.tagsWindow.SetHeight(m.height - 8 - m.status.Height())
	m.tagsWindow.SetWidth(m.width - 2)
}

func (m *TagListScreenModel) View() tea.View {
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
				m.drawTags(),
				m.status.View(),
			),
		)

	fmt.Fprint(sb, out)

	v.SetContent(sb.String())

	return v
}

func (m *TagListScreenModel) drawHeader() string {
	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		figlet.Figlet,
		" ",
		m.keysWindow.View(),
	)

	return out
}

func (m *TagListScreenModel) drawTags() string {
	tagLeftTitle := fmt.Sprintf("%q Tags", m.selectedRepositoryName)
	if m.tagListUI.IsFiltered() {
		tagLeftTitle += lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			Render(" /" + m.tagListUI.FilterValue())
	}
	m.tagsWindow.SetLeftTitle(tagLeftTitle)

	m.tagsWindow.SetRightTitle(
		lo.Ternary(
			m.isLoading,
			m.spinner.View(),
			ui.CountKind(
				len(m.tagList),
				"tag",
				"tags",
			),
		),
	)

	m.tagsWindow.SetContent(func(width, height int) string {
		if len(m.tagList) == 0 && m.isLoading {
			return "Loading..."
		}

		if len(m.tagList) == 0 {
			return "No tags found"
		}

		m.tagListUI.SetWidth(width)
		m.tagListUI.SetHeight(height)
		return m.tagListUI.View()
	})

	return m.tagsWindow.View()
}

func (m *TagListScreenModel) toggleMarkedTag() {
	if len(m.tagList) == 0 {
		return
	}

	targetTag := m.tagList[m.tagListUI.GlobalIndex()]

	if slices.ContainsFunc(*m.markedTags, func(item *Tag) bool { return item.Name == targetTag.Name }) {
		*m.markedTags = lo.Filter(*m.markedTags, func(item *Tag, _ int) bool {
			return item.Name != targetTag.Name
		})
	} else {
		*m.markedTags = append(*m.markedTags, targetTag)
	}
}
