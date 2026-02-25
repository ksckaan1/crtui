package ui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/samber/lo"
)

type TabbedWindowConfig struct {
	Width, Height  int
	ActiveTabIndex int
	Tabs           []string
	Content        string
	BorderColor    color.Color
	ActiveTabColor color.Color
}

type TabbedWindow struct {
	width, height       int
	activeTabIndex      int
	tabs                []string
	borderColor         color.Color
	activeTabTitleColor color.Color
	content             string
	vp                  viewport.Model
}

func NewTabbedWindow() *TabbedWindow {
	return &TabbedWindow{
		width:               0,
		height:              0,
		activeTabIndex:      0,
		tabs:                []string{},
		borderColor:         lipgloss.Color("#383838"),
		activeTabTitleColor: lipgloss.Color("#FFFFFF"),
		content:             "",
		vp:                  viewport.New(),
	}
}

func (t *TabbedWindow) Init() tea.Cmd {
	return nil
}

func (t *TabbedWindow) Update(msg tea.Msg) (*TabbedWindow, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "left" && len(t.tabs) != 0:
			t.activeTabIndex--
			if t.activeTabIndex < 0 {
				t.activeTabIndex = len(t.tabs) - 1
			}

		case msg.String() == "right" && len(t.tabs) != 0:
			t.activeTabIndex++
			if t.activeTabIndex >= len(t.tabs) {
				t.activeTabIndex = 0
			}
		}
	}

	var vpCmd tea.Cmd
	t.vp, vpCmd = t.vp.Update(msg)

	cmds := []tea.Cmd{vpCmd}

	return t, tea.Batch(cmds...)
}

func (t *TabbedWindow) View() string {
	if len(t.tabs) == 0 {
		return "No tabs found"
	}

	left := 0

	tabRanges := make([][2]int, len(t.tabs))
	last := -1

	drawnTabs := []string{}
	for i, title := range t.tabs {
		tab := t.drawTab(title, i == t.activeTabIndex, t.borderColor, t.activeTabTitleColor)
		drawnTabs = append(drawnTabs, tab)
		tabRanges[i] = [2]int{last + 1, last + lipgloss.Width(tab)}
		last += lipgloss.Width(tab)
	}

	activeTabRange := tabRanges[t.activeTabIndex]

	if activeTabRange[1] > left+t.width-2 {
		left = max(activeTabRange[1]-(t.width-6), 0)
	}

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		drawnTabs...,
	)

	outLen := lipgloss.Width(out)
	overflow := left+t.width-2 < outLen

	pd := 0
	if left > 0 {
		pd++
	}

	if overflow {
		pd++
	}

	out = t.cropHorizontal(out, left, t.width-pd-2)

	diff := max(t.width-outLen-2, 0)
	placeholder := make([]string, diff)
	borderText := lipgloss.NewStyle().Foreground(t.borderColor)

	for i := range diff {
		placeholder[i] = "\n\n─"
	}

	topBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		borderText.Render("\n\n╭"),
		lo.Ternary(left > 0, borderText.Render("«\n«\n«"), ""),
		out,
		lo.Ternary(overflow, borderText.Render("»\n»\n»"), ""),
		borderText.Render(lipgloss.JoinHorizontal(lipgloss.Top, placeholder...)),
		borderText.Render("\n\n╮"),
	)

	border := lipgloss.RoundedBorder()
	// BODY
	bodyStyle := lipgloss.NewStyle().
		Border(border, false, true, true, true).
		BorderForeground(t.borderColor).
		Width(t.width).
		Height(t.height - 1)

	t.vp.SetWidth(t.width - 2)
	t.vp.SetHeight(t.height - 2)
	t.vp.SetContent(t.content)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		bodyStyle.Render(t.vp.View()),
	)
}

func (t *TabbedWindow) SetContent(cb func(width, height, activeTabIndex int) string) {
	t.content = cb(t.width-2, t.height-4, t.activeTabIndex)
}

func (t *TabbedWindow) Width() int {
	return t.width
}

func (t *TabbedWindow) SetWidth(width int) {
	t.width = width
}

func (t *TabbedWindow) Height() int {
	return t.height
}

func (t *TabbedWindow) SetHeight(height int) {
	t.height = height
}

func (t *TabbedWindow) ActiveTabIndex() int {
	return t.activeTabIndex
}

func (t *TabbedWindow) SetActiveTabIndex(index int) {
	t.activeTabIndex = index
}

func (t *TabbedWindow) Tabs() []string {
	return t.tabs
}

func (t *TabbedWindow) SetTabs(tabs []string) {
	t.tabs = tabs
}

func (t *TabbedWindow) drawTab(title string, active bool, borderColor, activeTabColor color.Color) string {
	pLine := strings.Repeat("─", len(title))
	pEmpty := strings.Repeat(" ", len(title))

	borderText := lipgloss.NewStyle().
		Foreground(borderColor)

	top := borderText.Render(fmt.Sprintf("╭─%s─╮", pLine))

	middle := borderText.Render("│ ") +
		lo.Ternary(
			active,
			lipgloss.NewStyle().
				Foreground(activeTabColor).
				Render(title),
			borderText.Render(title),
		) +
		borderText.Render(" │")

	bottom := ""

	if active {
		bottom = borderText.Render(fmt.Sprintf("╯ %s ╰", pEmpty))
	} else {
		bottom = borderText.Render(fmt.Sprintf("┴%s──┴", pLine))
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		top,
		middle,
		bottom,
	)
}

func (t *TabbedWindow) cropHorizontal(s string, left, length int) string {
	if length == 0 {
		return ""
	}

	width := min(lipgloss.Width(s), length)

	vp := viewport.New(
		viewport.WithWidth(width),
		viewport.WithHeight(lipgloss.Height(s)),
	)

	vp.SetContent(s)
	vp.SetXOffset(left)

	return vp.View()
}
