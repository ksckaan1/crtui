package taglist

import (
	"context"
	"errors"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/internal/core/customerrors"
)

type tagListResult struct {
	tagList []string
	err     error
}

func (m *TagListScreenModel) fetchTagList() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		tagList, err := m.rc.ListTags(ctx, m.selectedRepositoryName)

		return tagListResult{
			tagList: tagList.Tags,
			err:     err,
		}
	}
}

type deleteTagsResult struct {
	err error
}

func (m *DeleteTagsPopup) deleteTags() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		for _, tagName := range m.tagNames {
			err := m.rc.DeleteTag(ctx, m.repositoryName, tagName)
			if err != nil {
				errStatusCode, ok := errors.AsType[customerrors.ErrStatusCode](err)
				if !ok || errStatusCode.StatusCode != http.StatusNotFound {
					return deleteTagsResult{
						err: err,
					}
				}
			}
		}

		return deleteTagsResult{}
	}
}
