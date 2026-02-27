package registrylist

import "charm.land/bubbles/v2/key"

type keyMapT struct {
	Navigate       key.Binding
	SelectRegistry key.Binding
	Refresh        key.Binding
	Filter         key.Binding
	New            key.Binding
	Esc            key.Binding
	Quit           key.Binding
}

var keyMap = keyMapT{
	Navigate: key.NewBinding(
		key.WithHelp("↑ ↓", "Navigate"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Refresh"),
	),
	SelectRegistry: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}
