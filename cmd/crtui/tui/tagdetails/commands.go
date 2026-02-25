package tagdetails

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/internal/core/models"
)

type tagMsg struct {
	tag *models.Tag
	err error
}

func (m *Model) fetchTag() tea.Cmd {
	return func() tea.Msg {
		tag, err := m.rc.GetTag(context.Background(), m.repositoryName, m.tagName)
		return tagMsg{tag: tag, err: err}
	}
}
