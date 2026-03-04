package figlet

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/version"
)

var Figlet = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Faint(false).Render(FigletText)

// const FigletText = `  ____ ____ _____ _   _ ___
// / ____|  _ \_   _| | | |_ _|
// | |   | |_) || | | | | || |
// | |___|  _ < | | | |_| || |
// \_____|_| \_\|_|  \___/|___|`

const figletText = ` %s╭─╮
╭────╮╭────╮╭─╮  ╭─╮╭─╮├─┤
│ ╭──╯│ ╭─┬┴╯ ╰─╮│ ││ ││ │
│ │   │ │ ╰─╮ ╭─╯│ ││ ││ │
│ ╰──╮│ │   │ ╰─╮│ ╰╯ ││ │
╰────╯╰─╯   ╰───╯╰────╯╰─╯`

var FigletText = func() string {
	spaceWidth := 22
	versionWidth := len(version.Version)
	padding := spaceWidth - versionWidth

	return fmt.Sprintf(figletText, version.Version+strings.Repeat(" ", padding))
}()
