package dockerhubclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"resty.dev/v3"
)

const apiBaseURL = "https://hub.docker.com"

// Docker Hub caps anonymous pagination at an offset of roughly 100 results.
const (
	pageSizeAuthenticated = "100"
	pageSizeAnonymous     = "25"
)

type Client struct {
	client   *resty.Client
	username string
	password string
}

// Repository represents a repository (image) on Docker Hub.
type Repository struct {
	Name      string
	Namespace string
	IsPrivate bool
}

type dockerHubRepository struct {
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	RepositoryType string `json:"repository_type"`
	IsPrivate      bool   `json:"is_private"`
}

type listResponse struct {
	Results []dockerHubRepository `json:"results"`
	Next    *string               `json:"next"`
}

func New(username, password string) *Client {
	return NewWithBaseURL(username, password, apiBaseURL)
}

func NewWithBaseURL(username, password, baseURL string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetDisableWarn(true)

	if username != "" && password != "" {
		client = client.SetBasicAuth(username, password)
	}

	return &Client{
		client:   client,
		username: username,
		password: password,
	}
}

// ListRepositories returns the image repositories of the configured namespace
// as Docker Hub repository names (namespace/name). Docker Hub does not expose
// a catalog endpoint, so this uses the public hub.docker.com API.
func (c *Client) ListRepositories(ctx context.Context) ([]Repository, error) {
	path := fmt.Sprintf("/v2/repositories/%s/", url.PathEscape(c.username))

	pageSize := pageSizeAnonymous
	if c.password != "" {
		pageSize = pageSizeAuthenticated
	}

	repos := make([]Repository, 0)

	next := path
	for next != "" {
		var page listResponse

		req := c.client.R().
			SetContext(ctx).
			SetResult(&page)

		if next == path {
			req = req.SetQueryParam("page_size", pageSize)
		}

		resp, err := req.Get(next)
		if err != nil {
			return nil, fmt.Errorf("request.Get: %w", err)
		}

		if resp.StatusCode() != http.StatusOK {
			body := string(resp.Bytes())
			resp.Body.Close()

			// Docker Hub refuses anonymous pagination past an offset of ~100
			// results. Stop gracefully so the public repositories fetched so
			// far are still returned.
			if c.password == "" &&
				resp.StatusCode() == http.StatusForbidden &&
				strings.Contains(body, "pagination offset too large") {
				break
			}

			return nil, fmt.Errorf("docker hub api returned %s", resp.Status())
		}

		for _, r := range page.Results {
			if r.Name == "" || (r.RepositoryType != "" && r.RepositoryType != "image") {
				continue
			}

			repos = append(repos, Repository{
				Name:      r.Name,
				Namespace: r.Namespace,
				IsPrivate: r.IsPrivate,
			})
		}

		resp.Body.Close()

		if page.Next != nil {
			next = *page.Next
		} else {
			next = ""
		}
	}

	return repos, nil
}

// Validate verifies that the configured namespace exists and that the
// credentials are accepted. It stops after the first page.
func (c *Client) Validate(ctx context.Context) error {
	var page listResponse

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&page).
		SetQueryParam("page_size", "1").
		Get(fmt.Sprintf("/v2/repositories/%s/", url.PathEscape(c.username)))
	if err != nil {
		return fmt.Errorf("request.Get: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("docker hub api returned %s", resp.Status())
	}

	return nil
}
