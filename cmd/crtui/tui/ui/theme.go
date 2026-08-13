package ui

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
)

// lightDark resolves a (light, dark) color pair to the appropriate color
// based on whether the terminal background is light or dark. It is recreated
// by SetDarkBackground whenever the background color changes.
var lightDark lipgloss.LightDarkFunc

// Semantic colors used across the UI. Each one is resolved from a light/dark
// pair through lipgloss.LightDark.
var (
	PrimaryColor        color.Color
	SecondaryColor      color.Color
	BorderColor         color.Color
	TextColor           color.Color
	SubtitleColor       color.Color
	TitleColor          color.Color
	ActiveTabTitleColor color.Color
	KeyColor            color.Color
	DescriptionColor    color.Color
	FaintColor          color.Color
	PopupTitleColor     color.Color

	// ThemeVersion is incremented whenever the theme changes. Cached render
	// outputs (e.g. KeysWindow content) can check it to rebuild themselves.
	ThemeVersion = 0
)

type themePair struct {
	light, dark color.Color
}

func (p themePair) resolve() color.Color {
	return lightDark(p.light, p.dark)
}

var (
	primaryPair = themePair{
		light: lipgloss.Color("#C2185B"),
		dark:  lipgloss.Color("#E60076"),
	}
	secondaryPair = themePair{
		light: lipgloss.Color("#8A8A8A"),
		dark:  lipgloss.Color("#383838"),
	}
	borderPair = themePair{
		light: lipgloss.Color("#9E9E9E"),
		dark:  lipgloss.Color("#383838"),
	}
	textPair = themePair{
		light: lipgloss.Color("#111111"),
		dark:  lipgloss.Color("#FFFFFF"),
	}
	subtitlePair = themePair{
		light: lipgloss.Color("#444444"),
		dark:  lipgloss.Color("#EAEAEA"),
	}
	titlePair = themePair{
		light: lipgloss.Color("#6B6B6B"),
		dark:  lipgloss.Color("#959595"),
	}
	activeTabTitlePair = themePair{
		light: lipgloss.Color("#000000"),
		dark:  lipgloss.Color("#FFFFFF"),
	}
	keyPair = themePair{
		light: lipgloss.Color("#000000"),
		dark:  lipgloss.Color("#FFFFFF"),
	}
	descriptionPair = themePair{
		light: lipgloss.Color("#555555"),
		dark:  lipgloss.Color("#626262"),
	}
	faintPair = themePair{
		light: lipgloss.Color("#9E9E9E"),
		dark:  lipgloss.Color("#444444"),
	}
	popupTitlePair = themePair{
		light: lipgloss.Color("#000000"),
		dark:  lipgloss.Color("#FFFFFF"),
	}
)

// SetDarkBackground resolves all semantic colors for the given terminal
// background and rebuilds anything that caches theme-dependent output.
func SetDarkBackground(isDark bool) {
	lightDark = lipgloss.LightDark(isDark)

	PrimaryColor = primaryPair.resolve()
	SecondaryColor = secondaryPair.resolve()
	BorderColor = borderPair.resolve()
	TextColor = textPair.resolve()
	SubtitleColor = subtitlePair.resolve()
	TitleColor = titlePair.resolve()
	ActiveTabTitleColor = activeTabTitlePair.resolve()
	KeyColor = keyPair.resolve()
	DescriptionColor = descriptionPair.resolve()
	FaintColor = faintPair.resolve()
	PopupTitleColor = popupTitlePair.resolve()

	LoadingSpinner = newLoadingSpinner()

	ThemeVersion++
}

func init() {
	SetDarkBackground(lipgloss.HasDarkBackground(os.Stdin, os.Stdout))
}
