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

func (m *TagDetailsScreenModel) fetchTag() tea.Cmd {
	return func() tea.Msg {
		tag, err := m.rc.GetTag(context.Background(), m.repositoryName, m.tagName)
		return tagMsg{tag: tag, err: err}
	}
}

type deleteTagMsg struct {
	err error
}

func (m *DeleteTagPopupModel) deleteTag() tea.Cmd {
	return func() tea.Msg {
		err := m.rc.DeleteTag(context.Background(), m.repositoryName, m.tagName)
		return deleteTagMsg{err: err}
	}
}
