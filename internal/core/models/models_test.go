package models

import (
	"testing"
	"time"

	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	registry := &Registry{
		BaseURL:       "https://registry.example.com",
		SupportsHTTP3: true,
		Username:      "user",
		Status:        registrystatus.Online,
		Took:          100 * time.Millisecond,
	}

	require.Equal(t, "https://registry.example.com", registry.BaseURL)
	require.True(t, registry.SupportsHTTP3)
	require.Equal(t, "user", registry.Username)
	require.Equal(t, registrystatus.Online, registry.Status)
	require.Equal(t, 100*time.Millisecond, registry.Took)
}

func TestTagList(t *testing.T) {
	tagList := &TagList{
		Name: "test/repo",
		Tags: []string{"v1.0.0", "v2.0.0", "latest"},
	}

	require.Equal(t, "test/repo", tagList.Name)
	require.Len(t, tagList.Tags, 3)
}

func TestTag(t *testing.T) {
	tag := &Tag{
		RepositoryName: "test/repo",
		TagName:        "v1.0.0",
		Platforms: []Platform{
			{
				Config: Config{
					Architecture: "amd64",
					OS:           "linux",
				},
			},
			{
				Config: Config{
					Architecture: "arm64",
					OS:           "linux",
				},
			},
		},
	}

	require.Equal(t, "test/repo", tag.RepositoryName)
	require.Equal(t, "v1.0.0", tag.TagName)
	require.Len(t, tag.Platforms, 2)
}

func TestPlatform(t *testing.T) {
	platform := Platform{
		Annotations: map[string]string{
			"org.opencontainers.image.source": "https://github.com/test/repo",
		},
	}

	require.NotNil(t, platform.Annotations)
	require.Contains(t, platform.Annotations, "org.opencontainers.image.source")
}

func TestConfig(t *testing.T) {
	config := Config{
		Digest:       "sha256:abc123",
		Architecture: "amd64",
		OS:           "linux",
		Author:       "Test Author",
		Created:      time.Now(),
		Entrypoint:   []string{"/bin/app"},
		Cmd:          []string{"--flag"},
		Env:          []string{"VAR=value"},
		WorkingDir:   "/app",
		Labels:       map[string]string{"version": "1.0"},
		User:         "root",
	}

	require.Equal(t, "sha256:abc123", config.Digest)
	require.Equal(t, "amd64", config.Architecture)
	require.Equal(t, "linux", config.OS)
	require.Equal(t, "Test Author", config.Author)
	require.Equal(t, []string{"/bin/app"}, config.Entrypoint)
	require.Equal(t, []string{"--flag"}, config.Cmd)
	require.Equal(t, []string{"VAR=value"}, config.Env)
	require.Equal(t, "/app", config.WorkingDir)
	require.Equal(t, "root", config.User)
}

func TestHistory(t *testing.T) {
	history := History{
		Author:     "Build System",
		Created:    time.Now(),
		CreatedBy:  "go build",
		Comment:    "Initial build",
		EmptyLayer: false,
	}

	require.Equal(t, "Build System", history.Author)
	require.Equal(t, "go build", history.CreatedBy)
	require.False(t, history.EmptyLayer)
}

func TestLayer(t *testing.T) {
	layer := Layer{
		MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
		Size:      1024000,
		Digest:    "sha256:def456",
	}

	require.Equal(t, "application/vnd.docker.image.rootfs.diff.tar.gzip", layer.MediaType)
	require.Equal(t, 1024000, layer.Size)
	require.Equal(t, "sha256:def456", layer.Digest)
}

func TestRootFS(t *testing.T) {
	rootFS := RootFS{
		Type:    "layers",
		DiffIDs: []string{"sha256:abc123", "sha256:def456"},
	}

	require.Equal(t, "layers", rootFS.Type)
	require.Len(t, rootFS.DiffIDs, 2)
}
