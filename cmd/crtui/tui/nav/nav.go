package nav

import tea "github.com/charmbracelet/bubbletea"

func Navigate(m tea.Model) (tea.Model, tea.Cmd) {
	return m, m.Init()
}
