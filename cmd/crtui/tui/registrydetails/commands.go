package registrydetails

import (
	"context"

	tea "charm.land/bubbletea/v2"
)

type repositoryListResult struct {
	repositoryList []string
	err            error
}

func (m *RegistryDetailsScreenModel) fetchRepositoryList() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		repositoryList, err := m.rc.ListRepositories(ctx)

		return repositoryListResult{
			repositoryList: repositoryList,
			err:            err,
		}
	}
}

type tagListResult struct {
	tagList []string
	err     error
}

func (m *RegistryDetailsScreenModel) fetchTagList() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		tagList, err := m.rc.ListTags(ctx, *m.selectedRepository)

		return tagListResult{
			tagList: tagList.Tags,
			err:     err,
		}
	}
}
