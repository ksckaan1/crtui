package registryclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/samber/lo"
	"resty.dev/v3"
)

type RegistryClient struct {
	client           *resty.Client
	baseURL          string
	username         string
	password         string
	http3Transport   *http3.Transport
	defaultTransport *http.Transport
	supportsHTTP3    bool
}

func New(baseURL, username, password string) *RegistryClient {
	client := resty.New().
		SetBasicAuth(username, password).
		SetBaseURL(baseURL)

	return &RegistryClient{
		client:   client,
		baseURL:  baseURL,
		username: username,
		password: password,
		defaultTransport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
		http3Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{},
			QUICConfig: &quic.Config{
				MaxIdleTimeout:  45 * time.Second,
				KeepAlivePeriod: 30 * time.Second,
			},
		},
	}
}

func (r *RegistryClient) GetRegistryInfo(ctx context.Context) (*models.Registry, error) {
	startTime := time.Now()
	resp, err := r.client.R().SetContext(ctx).
		Get("/v2/")
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}
	defer resp.Body.Close()

	supportsHTTP3 := strings.Contains(resp.Header().Get("alt-svc"), "h3=")
	r.setTransport(supportsHTTP3)

	authSuccess := resp.StatusCode() == http.StatusOK
	took := time.Since(startTime)

	return &models.Registry{
		BaseURL:       r.baseURL,
		SupportsHTTP3: supportsHTTP3,
		Username:      r.username,
		AuthSuccess:   authSuccess,
		Took:          took,
	}, nil
}

func (r *RegistryClient) setTransport(supportsHTTP3 bool) {
	if supportsHTTP3 && !r.supportsHTTP3 {
		r.supportsHTTP3 = supportsHTTP3
		r.client = r.client.SetTransport(r.http3Transport)
		return
	}

	if !supportsHTTP3 && r.supportsHTTP3 {
		r.supportsHTTP3 = supportsHTTP3
		r.client = r.client.SetTransport(r.defaultTransport)
	}
}

func (r *RegistryClient) Close() error {
	err := r.http3Transport.Close()
	if err != nil {
		return fmt.Errorf("http3Transport.Close: %w", err)
	}

	err = r.client.Close()
	if err != nil {
		return fmt.Errorf("client.Close: %w", err)
	}

	return nil
}

func (r *RegistryClient) ListRepositories(ctx context.Context) ([]string, error) {
	var result struct {
		Repositories []string `json:"repositories"`
	}

	resp, err := r.client.R().SetContext(ctx).
		SetResult(&result).
		Get("/v2/_catalog")
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return result.Repositories, nil
}

func (r *RegistryClient) ListTags(ctx context.Context, repository string) (*models.TagList, error) {
	var result struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}

	uri := path.Join("v2", repository, "tags/list")

	resp, err := r.client.R().SetContext(ctx).
		SetResult(&result).
		Get(uri)
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return &models.TagList{
		Name: result.Name,
		Tags: result.Tags,
	}, nil
}

func (r *RegistryClient) GetTag(ctx context.Context, repositoryName, tagName string) (*models.Tag, error) {
	uri := path.Join("v2", repositoryName, "manifests", tagName)

	resp, err := r.client.R().SetContext(ctx).
		SetHeader("Accept", "application/vnd.oci.image.index.v1+json, application/vnd.oci.image.manifest.v1+json, application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json").
		Get(uri)
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	switch resp.Header().Get("Content-Type") {
	case "application/vnd.oci.image.manifest.v1+json":
		return r.processImageManifest(ctx, repositoryName, tagName, resp.Bytes())
	case "application/vnd.oci.image.index.v1+json":
		return r.processImageIndex(ctx, repositoryName, tagName, resp.Bytes())
	default:
		return nil, fmt.Errorf("unexpected content type: %s", resp.Header().Get("Content-Type"))
	}
}

func (r *RegistryClient) processImageManifest(ctx context.Context, repositoryName, tagName string, body []byte) (*models.Tag, error) {
	var parsedBody imageManifestV1

	err := json.Unmarshal(body, &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	configBlob, err := r.getConfigBlob(ctx, repositoryName, parsedBody.Config.Digest)
	if err != nil {
		return nil, fmt.Errorf("getConfigBlob: %w", err)
	}

	tag := &models.Tag{
		RepositoryName: repositoryName,
		TagName:        tagName,
		Platforms: []models.Platform{
			{
				Config: models.Config{
					Digest:       parsedBody.Config.Digest,
					Architecture: configBlob.Architecture,
					OS:           configBlob.OS,
					Author:       configBlob.Author,
					Created:      configBlob.Created,
					History: lo.Map(configBlob.History, func(item blobHistory, _ int) models.History {
						return models.History{
							Created:   item.Created,
							CreatedBy: item.CreatedBy,
							Comment:   item.Comment,
							Author:    item.Author,
						}
					}),
					Entrypoint: configBlob.Config.Entrypoint,
					Env:        configBlob.Config.Env,
					Labels:     configBlob.Config.Labels,
					User:       configBlob.Config.User,
				},
				Layers: lo.Map(parsedBody.Layers, func(item imageManifestV1Layer, _ int) models.Layer {
					return models.Layer{
						MediaType: item.MediaType,
						Size:      item.Size,
						Digest:    item.Digest,
					}
				}),
				Annotations: parsedBody.Annotations,
			},
		},
	}

	return tag, nil
}

func (r *RegistryClient) processImageIndex(ctx context.Context, repositoryName, tagName string, body []byte) (*models.Tag, error) {
	var parsedBody imageIndexV1

	err := json.Unmarshal(body, &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	tag := &models.Tag{
		RepositoryName: repositoryName,
		TagName:        tagName,
	}

	for _, manifest := range parsedBody.Manifests {
		t, err := r.getManifestByDigest(ctx, repositoryName, manifest.Digest)
		if err != nil {
			return nil, fmt.Errorf("getManifestByDigest: %w", err)
		}

		platform := t.Platforms[0]

		tag.Platforms = append(tag.Platforms, models.Platform{
			Config:      platform.Config,
			Layers:      platform.Layers,
			Annotations: platform.Annotations,
		})
	}

	return tag, nil
}

func (r *RegistryClient) getConfigBlob(ctx context.Context, repositoryName, digest string) (*blob, error) {
	uri := path.Join("v2", repositoryName, "blobs", digest)

	resp, err := r.client.R().
		SetContext(ctx).
		Get(uri)
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var result blob

	err = json.Unmarshal(resp.Bytes(), &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &result, nil
}

func (r *RegistryClient) getManifestByDigest(ctx context.Context, repositoryName, digest string) (*models.Tag, error) {
	uri := path.Join("v2", repositoryName, "manifests", digest)

	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/vnd.oci.image.index.v1+json, application/vnd.docker.distribution.manifest.list.v2+json, application/vnd.docker.distribution.manifest.v2+json, application/vnd.oci.image.manifest.v1+json").
		Get(uri)
	if err != nil {
		return nil, fmt.Errorf("request.Get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var parsedBody imageManifestV1

	err = json.Unmarshal(resp.Bytes(), &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	configBlob, err := r.getConfigBlob(ctx, repositoryName, parsedBody.Config.Digest)
	if err != nil {
		return nil, fmt.Errorf("getConfigBlob: %w", err)
	}

	tag := &models.Tag{
		Platforms: []models.Platform{
			{
				Config: models.Config{
					Digest:       parsedBody.Config.Digest,
					Architecture: configBlob.Architecture,
					OS:           configBlob.OS,
					Author:       configBlob.Author,
					Created:      configBlob.Created,
					History: lo.Map(configBlob.History, func(item blobHistory, _ int) models.History {
						return models.History{
							Created:   item.Created,
							CreatedBy: item.CreatedBy,
							Comment:   item.Comment,
							Author:    item.Author,
						}
					}),
					Entrypoint: configBlob.Config.Entrypoint,
					Env:        configBlob.Config.Env,
					Labels:     configBlob.Config.Labels,
					User:       configBlob.Config.User,
				},
				Layers: lo.Map(parsedBody.Layers, func(item imageManifestV1Layer, _ int) models.Layer {
					return models.Layer{
						MediaType: item.MediaType,
						Size:      item.Size,
						Digest:    item.Digest,
					}
				}),
				Annotations: parsedBody.Annotations,
			},
		},
	}

	return tag, nil
}

func (r *RegistryClient) DeleteTag(ctx context.Context, repositoryName, tagName string) error {
	headURI := path.Join("v2", repositoryName, "manifests", tagName)

	headResp, err := r.client.R().SetContext(ctx).
		SetHeader("Accept", "application/vnd.oci.image.index.v1+json, application/vnd.oci.image.manifest.v1+json, application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json").
		Head(headURI)
	if err != nil {
		return fmt.Errorf("request.Head: %w", err)
	}
	defer headResp.Body.Close()

	if headResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", headResp.StatusCode())
	}

	dockerContentDigest := headResp.Header().Get("Docker-Content-Digest")

	deleteURI := path.Join("v2", repositoryName, "manifests", dockerContentDigest)
	deleteResp, err := r.client.R().SetContext(ctx).
		SetHeader("Accept", "application/vnd.oci.image.index.v1+json, application/vnd.oci.image.manifest.v1+json, application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json").
		Delete(deleteURI)
	if err != nil {
		return fmt.Errorf("request.Head: %w", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode() != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", deleteResp.StatusCode())
	}

	return nil
}
