package githubclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"resty.dev/v3"
)

const apiBaseURL = "https://api.github.com"

type Client struct {
	client *resty.Client
	owner  string
	token  string
}

// Package represents a container package (repository) on ghcr.io.
type Package struct {
	Name       string
	Visibility string
}

type githubPackage struct {
	Name        string `json:"name"`
	PackageType string `json:"package_type"`
	Visibility  string `json:"visibility"`
}

func New(owner, token string) *Client {
	return NewWithBaseURL(owner, token, apiBaseURL)
}

func NewWithBaseURL(owner, token, baseURL string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetDisableWarn(true)

	return &Client{
		client: client,
		owner:  owner,
		token:  token,
	}
}

// ListContainerPackages returns every container package owned by the configured
// owner as a ghcr.io repository name (owner/name) together with its visibility.
// Listing packages requires a GitHub token; public packages are returned as well.
func (c *Client) ListContainerPackages(ctx context.Context) ([]Package, error) {
	isOrg, err := c.isOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("detect owner type: %w", err)
	}

	path := fmt.Sprintf("/users/%s/packages", url.PathEscape(c.owner))
	if isOrg {
		path = fmt.Sprintf("/orgs/%s/packages", url.PathEscape(c.owner))
	}

	repos := make([]Package, 0)

	next := path
	for next != "" {
		var page []githubPackage

		resp, err := c.newRequest(ctx).
			SetResult(&page).
			SetQueryParam("package_type", "container").
			SetQueryParam("per_page", "100").
			Get(next)
		if err != nil {
			return nil, fmt.Errorf("request.Get: %w", err)
		}

		if resp.StatusCode() != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("github api returned %s", resp.Status())
		}

		for _, p := range page {
			if p.PackageType != "" && p.PackageType != "container" {
				continue
			}
			if p.Name != "" {
				repos = append(repos, Package{
					Name:       fmt.Sprintf("%s/%s", c.owner, p.Name),
					Visibility: p.Visibility,
				})
			}
		}

		resp.Body.Close()

		next = nextLink(resp.Header().Get("Link"))
	}

	return repos, nil
}

// Validate verifies that the configured owner exists and that the token can be
// used to list container packages. It stops after the first page.
func (c *Client) Validate(ctx context.Context) error {
	isOrg, err := c.isOrg(ctx)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%s/packages", url.PathEscape(c.owner))
	if isOrg {
		path = fmt.Sprintf("/orgs/%s/packages", url.PathEscape(c.owner))
	}

	var page []githubPackage

	resp, err := c.newRequest(ctx).
		SetResult(&page).
		SetQueryParam("package_type", "container").
		SetQueryParam("per_page", "1").
		Get(path)
	if err != nil {
		return fmt.Errorf("request.Get: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("github api returned %s", resp.Status())
	}

	return nil
}

func (c *Client) isOrg(ctx context.Context) (bool, error) {
	resp, err := c.newRequest(ctx).
		Get(fmt.Sprintf("/orgs/%s", url.PathEscape(c.owner)))
	if err != nil {
		return false, fmt.Errorf("request.Get: %w", err)
	}
	resp.Body.Close()

	switch resp.StatusCode() {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("github api returned %s", resp.Status())
	}
}

func (c *Client) newRequest(ctx context.Context) *resty.Request {
	req := c.client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		SetHeader("User-Agent", "crtui")

	if c.token != "" {
		req.SetHeader("Authorization", "Bearer "+c.token)
	}

	return req
}

// nextLink returns the URL of the next page from a GitHub Link header, or an
// empty string when there is no next page.
func nextLink(linkHeader string) string {
	for _, part := range strings.Split(linkHeader, ",") {
		part = strings.TrimSpace(part)

		if !strings.Contains(part, `rel="next"`) {
			continue
		}

		if _, rest, ok := strings.Cut(part, "<"); ok {
			if next, _, ok := strings.Cut(rest, ">"); ok {
				return strings.TrimSpace(next)
			}
		}
	}

	return ""
}
