package figlet

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/ui"
	"github.com/ksckaan1/crtui/version"
)

var Figlet = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Faint(false).Render(FigletText)

const figletText = ` %s╭─╮
╭────╮╭────╮╭─╮  ╭─╮╭─╮├─┤
│ ╭──╯│ ╭─┬┴╯ ╰─╮│ ││ ││ │
│ │   │ │ ╰─╮ ╭─╯│ ││ ││ │
│ ╰──╮│ │   │ ╰─╮│ ╰╯ ││ │
╰────╯╰─╯   ╰───╯╰────╯╰─╯`

var FigletText = func() string {
	spaceWidth := 22
	versionText := version.Version
	versionWidth := len(versionText)

	if versionWidth > spaceWidth {
		versionText = versionText[:spaceWidth]
	}

	padding := spaceWidth - versionWidth

	return fmt.Sprintf(figletText, versionText+strings.Repeat(" ", padding))
}()
