package config

import (
	"fmt"
	"strings"

	"github.com/docker/cli/cli/config"
)

func listAutoDetectedAuths() ([]*Auth, error) {
	dockerCFG, err := config.Load(config.Dir())
	if err != nil {
		return nil, fmt.Errorf("config.Load: %w (docker config)", err)
	}

	dockerCreds, err := dockerCFG.GetAllCredentials()
	if err != nil {
		return nil, fmt.Errorf("cfg.GetAllCredentials: %w (docker config)", err)
	}

	dockerAuths := make([]*Auth, 0)

	for key, authConfig := range dockerCreds {
		if strings.HasPrefix(key, "https://index.docker.io/v1/") {
			continue
		}

		dockerAuths = append(dockerAuths, &Auth{
			URL:          key,
			Username:     authConfig.Username,
			Password:     authConfig.Password,
			AutoDetected: true,
		})
	}

	return dockerAuths, nil
}
