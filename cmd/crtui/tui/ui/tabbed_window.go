package ui

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

type TabbedWindowConfig struct {
	Width, Height  int
	ActiveTabIndex int
	Tabs           []string
	Content        string
	BorderColor    lipgloss.Color
	ActiveTabColor lipgloss.Color
}

func NewTabbedWindow(cfg TabbedWindowConfig) string {
	borderColor := cmp.Or(cfg.BorderColor, lipgloss.Color("#383838"))
	activeTabColor := cmp.Or(cfg.ActiveTabColor, lipgloss.Color("#FFFFFF"))

	left := 0

	tabRanges := make([][2]int, len(cfg.Tabs))
	last := -1

	drawnTabs := []string{}
	for i, title := range cfg.Tabs {
		tab := drawTab(title, i == cfg.ActiveTabIndex, borderColor, activeTabColor)
		drawnTabs = append(drawnTabs, tab)
		tabRanges[i] = [2]int{last + 1, last + lipgloss.Width(tab)}
		last += lipgloss.Width(tab)
	}

	activeTabRange := tabRanges[cfg.ActiveTabIndex]

	if activeTabRange[1] > left+cfg.Width {
		left = max(activeTabRange[1]-(cfg.Width-4), 0)
	}

	out := lipgloss.JoinHorizontal(
		lipgloss.Top,
		drawnTabs...,
	)

	outLen := lipgloss.Width(out)
	overflow := left+cfg.Width < outLen

	pd := 0
	if left > 0 {
		pd++
	}

	if overflow {
		pd++
	}

	out = cropHorizontal(out, left, cfg.Width-pd)

	diff := max(cfg.Width-outLen, 0)
	placeholder := make([]string, diff)
	borderText := lipgloss.NewStyle().Foreground(borderColor)

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
		BorderForeground(borderColor).
		Width(cfg.Width).
		Height(cfg.Height - 1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		bodyStyle.Render(cfg.Content),
	)
}

func drawTab(title string, active bool, borderColor, activeTabColor lipgloss.Color) string {
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

func cropHorizontal(s string, left, length int) string {
	if length == 0 {
		return ""
	}

	width := min(lipgloss.Width(s), length)

	vp := viewport.New(width, lipgloss.Height(s))

	vp.SetContent(s)
	vp.SetXOffset(left)

	return vp.View()
}
