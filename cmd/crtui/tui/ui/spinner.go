package ui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"
)

var LoadingSpinner spinner.Spinner

func newLoadingSpinner() spinner.Spinner {
	const text = "loading"

	frames := make([]string, 0, len(text))

	for i := range text {
		var frame strings.Builder
		for j, r := range text {
			if j == i {
				frame.WriteString(lipgloss.NewStyle().Foreground(PrimaryColor).Render(string(r)))
			} else {
				frame.WriteString(lipgloss.NewStyle().Foreground(TextColor).Render(string(r)))
			}
		}
		frames = append(frames, frame.String())
	}

	return spinner.Spinner{
		Frames: frames,
		FPS:    time.Second / 7,
	}
}
