package registrytype

import (
	"net/url"
	"strings"
)

type RegistryType string

const (
	Docker    RegistryType = "docker"
	GitHub    RegistryType = "github"
	DockerHub RegistryType = "dockerhub"
)

// DockerHubRegistryHost is the canonical host that serves the Docker
// Distribution API for Docker Hub (index.docker.io/docker.io are aliases).
const DockerHubRegistryHost = "registry-1.docker.io"

func FromURL(rawURL string) RegistryType {
	u, err := url.Parse(rawURL)
	if err != nil {
		return Docker
	}

	if u.Scheme == "" {
		u, err = url.Parse("https://" + rawURL)
		if err != nil {
			return Docker
		}
	}

	switch strings.ToLower(u.Hostname()) {
	case "ghcr.io":
		return GitHub
	case "docker.io", "index.docker.io", "registry-1.docker.io", "registry.docker.io":
		return DockerHub
	default:
		return Docker
	}
}

// Normalize returns the canonical registry URL for the given registry URL.
// Docker Hub aliases (docker.io, index.docker.io, ...) are mapped to the
// host that serves the Docker Distribution API. Any other URL is returned
// unchanged.
func Normalize(rawURL string) string {
	if FromURL(rawURL) != DockerHub {
		return rawURL
	}

	return "https://" + DockerHubRegistryHost
}
