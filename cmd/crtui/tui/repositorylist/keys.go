package repositorylist

import "charm.land/bubbles/v2/key"

type keyMapT struct {
	Navigate key.Binding
	Select   key.Binding
	Refresh  key.Binding
	Filter   key.Binding
	Delete   key.Binding
	Esc      key.Binding
	Quit     key.Binding
}

var keyMap = keyMapT{
	Navigate: key.NewBinding(
		key.WithHelp("↑ ↓", "Navigate"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Refresh"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "Delete"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}
