package registrytype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromURL(t *testing.T) {
	cases := []struct {
		rawURL string
		want   RegistryType
	}{
		{"https://ghcr.io", GitHub},
		{"https://ghcr.io/", GitHub},
		{"http://ghcr.io", GitHub},
		{"GHCR.IO", GitHub},
		{"https://registry.example.com", Docker},
		{"http://localhost:5000", Docker},
		{"not a url", Docker},
		{"https://registry-1.docker.io", DockerHub},
		{"https://registry.docker.io", DockerHub},
		{"https://index.docker.io/v1/", DockerHub},
		{"https://index.docker.io", DockerHub},
		{"https://docker.io", DockerHub},
		{"DOCKER.IO", DockerHub},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, FromURL(tc.rawURL), tc.rawURL)
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		rawURL string
		want   string
	}{
		{"https://registry-1.docker.io", "https://registry-1.docker.io"},
		{"https://registry-1.docker.io/v1/", "https://registry-1.docker.io"},
		{"https://index.docker.io/v1/", "https://registry-1.docker.io"},
		{"https://index.docker.io", "https://registry-1.docker.io"},
		{"https://docker.io", "https://registry-1.docker.io"},
		{"index.docker.io", "https://registry-1.docker.io"},
		{"https://ghcr.io", "https://ghcr.io"},
		{"https://registry.example.com", "https://registry.example.com"},
		{"http://localhost:5000", "http://localhost:5000"},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, Normalize(tc.rawURL), tc.rawURL)
	}
}
