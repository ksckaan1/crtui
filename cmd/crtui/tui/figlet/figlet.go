package figlet

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
)

var Figlet = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Faint(false).Render(FigletText)

const FigletText = `  ____ ____ _____ _   _ ___
/ ____|  _ \_   _| | | |_ _|
| |   | |_) || | | | | || |
| |___|  _ < | | | |_| || |
\_____|_| \_\|_|  \___/|___|`
