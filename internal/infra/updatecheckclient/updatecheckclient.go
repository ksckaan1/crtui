package updatecheckclient

import (
	"context"
	"fmt"

	"github.com/ksckaan1/crtui/internal/core/customerrors"
	"github.com/ksckaan1/crtui/internal/core/models"
	"resty.dev/v3"
)

type UpdateCheckClient struct {
	client *resty.Client
}

func New() *UpdateCheckClient {
	c := resty.New().
		SetBaseURL("https://api.github.com").
		SetDisableWarn(true)

	return &UpdateCheckClient{
		client: c,
	}
}

func (c *UpdateCheckClient) GetLatestRelease(ctx context.Context, username, repo string) (*models.Release, error) {
	var result models.Release

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("User-Agent", "crtui").
		SetHeader("Accept", "application/vnd.github+json").
		SetResult(&result).
		SetPathParam("username", username).
		SetPathParam("repo", repo).
		Get("/repos/{username}/{repo}/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, customerrors.ErrStatusCode{
			StatusCode: resp.StatusCode(),
		}
	}

	return &result, nil
}
