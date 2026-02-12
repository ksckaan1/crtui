package figlet

import (
	"github.com/charmbracelet/lipgloss"
)

var Figlet = lipgloss.NewStyle().Faint(true).Render(FigletText)

const FigletText = `  ____ ____ _____ _   _ ___
/ ____|  _ \_   _| | | |_ _|
| |   | |_) || | | | | || |
| |___|  _ < | | | |_| || |
\_____|_| \_\|_|  \___/|___|`
