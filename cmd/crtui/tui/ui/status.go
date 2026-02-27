package ui

import (
	"fmt"
	"image/color"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

type StatusType string

const (
	OK      StatusType = "OK"
	Info    StatusType = "INF"
	Warning StatusType = "WRN"
	Error   StatusType = "ERR"
	Empty   StatusType = ""
)

type Status struct {
	badge      string
	statusText string
	tookStr    string
	updatedAt  string
	statusType StatusType
}

func NewStatus() *Status {
	return &Status{}
}

func (s *Status) SetStatus(statusType StatusType, statusText string, took ...time.Duration) tea.Cmd {
	tookDur := time.Duration(0)
	if len(took) > 0 {
		tookDur = took[0]
	}

	s.statusType = statusType

	if statusType == Empty {
		return tea.RequestWindowSize
	}

	var badgeColor color.Color

	switch statusType {
	case OK:
		badgeColor = lipgloss.Color("42")
	case Info:
		badgeColor = lipgloss.Color("#5dade2")
	case Warning:
		badgeColor = lipgloss.Color("226")
	case Error:
		badgeColor = lipgloss.Color("196")
	}

	s.badge = lipgloss.NewStyle().
		Foreground(badgeColor).
		MarginRight(1).
		Render(string(statusType))

	s.updatedAt = lipgloss.NewStyle().
		Faint(true).
		Bold(true).
		Render(time.Now().Format("15:04:05"))

	tookStr := ""
	if tookDur > 0 {
		tookStr = lipgloss.NewStyle().
			Padding(0, 1).
			Faint(true).
			Bold(true).
			Render(fmt.Sprintf("(%dms)", tookDur.Round(time.Millisecond).Milliseconds()))
	}

	s.tookStr = tookStr

	s.statusText = statusText

	return tea.RequestWindowSize
}

func (s *Status) View() string {
	if s.statusType == Empty {
		return ""
	}

	currentWidth, _, err := term.GetSize(os.Stdin.Fd())
	if err != nil {
		return "error when getting terminal size"
	}

	text := lipgloss.NewStyle().
		Width(
			currentWidth -
				lipgloss.Width(s.updatedAt) -
				lipgloss.Width(s.badge) -
				lipgloss.Width(s.tookStr) -
				2,
		).
		Render(s.statusText)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		s.badge,
		text,
		s.tookStr,
		s.updatedAt,
	)
}

func (s *Status) Height() int {
	v := s.View()
	if v == "" {
		return 0
	}
	return lipgloss.Height(v)
}
