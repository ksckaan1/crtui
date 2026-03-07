package tagdetails

import "charm.land/bubbles/v2/key"

type keyMapT struct {
	Scroll          key.Binding
	SwitchPlatforms key.Binding
	Refresh         key.Binding
	CopyPullCommand key.Binding
	DeleteTag       key.Binding
	Esc             key.Binding
	Quit            key.Binding
}

var keyMap = keyMapT{
	Scroll: key.NewBinding(
		key.WithHelp("↑ ↓", "Scroll"),
	),
	SwitchPlatforms: key.NewBinding(
		key.WithHelp("← →", "Switch Platforms"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Refresh"),
	),
	CopyPullCommand: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Copy Pull Command"),
	),
	DeleteTag: key.NewBinding(
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
