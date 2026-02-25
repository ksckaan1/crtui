package nav

import tea "charm.land/bubbletea/v2"

func Navigate(m tea.Model) (tea.Model, tea.Cmd) {
	return m, m.Init()
}
