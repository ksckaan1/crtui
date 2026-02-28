package registrylist

import "charm.land/bubbles/v2/key"

type keyMapT struct {
	Navigate       key.Binding
	SelectRegistry key.Binding
	Refresh        key.Binding
	Filter         key.Binding
	New            key.Binding
	Edit           key.Binding
	Delete         key.Binding
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
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "Delete"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}
