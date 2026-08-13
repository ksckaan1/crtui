package registrytype

import (
	"net/url"
	"strings"
)

type RegistryType string

const (
	Docker RegistryType = "docker"
	GitHub RegistryType = "github"
)

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

	if strings.EqualFold(u.Hostname(), "ghcr.io") {
		return GitHub
	}

	return Docker
}
