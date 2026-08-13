package repositorylist

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/internal/core/customerrors"
	"github.com/ksckaan1/crtui/internal/core/enums/registrytype"
	"github.com/ksckaan1/crtui/internal/infra/dockerhubclient"
	"github.com/ksckaan1/crtui/internal/infra/githubclient"
	"github.com/samber/lo"
)

type repositoryListResult struct {
	repositoryList []*Repository
	err            error
}

func (m *RepositoryListScreenModel) fetchRepositoryList() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		if m.registry.Type == registrytype.GitHub {
			ghClient := githubclient.New(m.registry.Username, m.registry.Password)

			packages, err := ghClient.ListContainerPackages(ctx)

			return repositoryListResult{
				repositoryList: lo.Map(packages, func(item githubclient.Package, _ int) *Repository {
					return &Repository{
						Name:       item.Name,
						Visibility: item.Visibility,
					}
				}),
				err: err,
			}
		}

		if m.registry.Type == registrytype.DockerHub {
			dhClient := dockerhubclient.New(m.registry.Username, m.registry.Password)

			repositories, err := dhClient.ListRepositories(ctx)

			return repositoryListResult{
				repositoryList: lo.Map(repositories, func(item dockerhubclient.Repository, _ int) *Repository {
					return &Repository{
						Name:       fmt.Sprintf("%s/%s", item.Namespace, item.Name),
						Visibility: lo.Ternary(item.IsPrivate, "private", "public"),
					}
				}),
				err: err,
			}
		}

		repositoryList, err := m.rc.ListRepositories(ctx)

		return repositoryListResult{
			repositoryList: lo.Map(repositoryList, func(item string, _ int) *Repository {
				return &Repository{Name: item}
			}),
			err: err,
		}
	}
}

type tagListResult struct {
	tagList []string
	err     error
}

func (m *DeleteTagsPopup) fetchTagList() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		tagList, err := m.rc.ListTags(ctx, m.repositoryName)

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
