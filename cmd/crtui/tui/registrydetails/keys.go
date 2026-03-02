package registrydetails

import "charm.land/bubbles/v2/key"

type keyMapT struct {
	Navigate   key.Binding
	SelectItem key.Binding
	Refresh    key.Binding
	Filter     key.Binding
	SwitchPane key.Binding
	Delete     key.Binding
	Esc        key.Binding
	Quit       key.Binding
}

var keyMap = keyMapT{
	Navigate: key.NewBinding(
		key.WithHelp("↑ ↓", "Navigate"),
	),
	SelectItem: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select item"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Refresh"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	SwitchPane: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Switch pane"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "Delete repo/tags"),
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
