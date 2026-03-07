package ui

import (
	"fmt"
	"math"
	"reflect"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

var (
	keysContainer = lipgloss.NewStyle().
			Width(100).
			PaddingLeft(1)
	keysKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(10)
	keysDescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262"))
)

type KeysWindow struct {
	keys    []key.Binding
	window  *Window
	content string
}

func NewKeysWindow() *KeysWindow {
	w := NewWindow()
	w.SetLeftTitle("Keys")

	return &KeysWindow{
		window: w,
	}
}

func (m *KeysWindow) SetWidth(width int) {
	if width == m.window.Width() {
		return
	}
	m.window.SetWidth(width)
	m.genContent()
}

func (m *KeysWindow) Width() int {
	return m.window.Width()
}

func (m *KeysWindow) SetHeight(height int) {
	if height == m.window.Height() {
		return
	}
	m.window.SetHeight(height)
	m.genContent()
}

func (m *KeysWindow) Height() int {
	return m.window.Height()
}

func (m *KeysWindow) SetKeyMap(keyMap any) {
	keys := []key.Binding{}

	v := reflect.ValueOf(keyMap)
	if v.Kind() != reflect.Struct {
		panic("key map must be struct")
	}

	for i := range v.Type().NumField() {
		fieldVal, ok := v.Field(i).Interface().(key.Binding)
		if !ok {
			continue
		}
		if fieldVal.Help().Desc == "" {
			continue
		}
		keys = append(keys, fieldVal)
	}

	m.keys = keys
	m.genContent()
}

func (m *KeysWindow) genContent() {
	m.window.SetContent(func(width, height int) string {
		cols := make([]string, 0, int(math.Ceil(float64(len(m.keys))/float64(height))))

		var col string
		for i, k := range m.keys {
			col += fmt.Sprintf("%s %s", keysKeyStyle.Render(k.Help().Key), keysDescriptionStyle.Render(k.Help().Desc))
			if i%height != 3 {
				col += "\n"
				continue
			}

			cols = append(cols, col, "  ")
			col = ""
		}

		if col != "" {
			cols = append(cols, col)
		}

		keysContent := lipgloss.JoinHorizontal(lipgloss.Top, cols...)

		return m.cropHorizontal(keysContainer.Render(keysContent), 0, width-1)
	})

	m.content = m.window.View()
}

func (m *KeysWindow) View() string {
	return m.content
}

func (t *KeysWindow) cropHorizontal(s string, left, length int) string {
	if length == 0 {
		return ""
	}

	width := min(lipgloss.Width(s), length)

	vp := viewport.New(
		viewport.WithWidth(width),
		viewport.WithHeight(lipgloss.Height(s)),
	)

	vp.SetContent(s)
	vp.SetXOffset(left)

	return vp.View()
}
