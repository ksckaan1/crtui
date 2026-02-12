package ui

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestKeysWindow(t *testing.T) {
	result := NewKeysWindow(KeysWindowConfig{
		Width:  40,
		Height: 4,
		Keys: []key.Binding{
			key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "for a")),
			key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "for b")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "for c")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "for d")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "for e")),
			key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "for f")),
			key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "for a")),
			key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "for b")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "for c")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "for d")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "for e")),
			key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "for f")),
		},
	})

	fmt.Println(result)
}
