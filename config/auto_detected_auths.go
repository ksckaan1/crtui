package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config"
)

func getDockerAuths() ([]*Auth, error) {
	_, err := os.Stat(config.Dir())
	if os.IsNotExist(err) {
		return nil, nil
	}

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

func getPodmanAuths() ([]*Auth, error) {
	userConfigDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("os.UserHomeDir: %w", err)
	}

	configFilePath := filepath.Join(userConfigDir, ".config", "containers", "auth.json")

	_, err = os.Stat(configFilePath)
	if err != nil {
		return nil, nil
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer configFile.Close()

	podmanCFG, err := config.LoadFromReader(configFile)
	if err != nil {
		return nil, fmt.Errorf("config.Load: %w (podman config)", err)
	}

	podmanCreds, err := podmanCFG.GetAllCredentials()
	if err != nil {
		return nil, fmt.Errorf("cfg.GetAllCredentials: %w (podman config)", err)
	}

	dockerAuths := make([]*Auth, 0)

	for key, authConfig := range podmanCreds {
		if strings.HasPrefix(key, "https://index.docker.io/v1/") {
			continue
		}

		if !(strings.HasPrefix(key, "https://") || strings.HasPrefix(key, "http://")) {
			key = "https://" + key
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

func listAutoDetectedAuths() ([]*Auth, error) {
	dockerAuths, err := getDockerAuths()
	if err != nil {
		return nil, fmt.Errorf("getDockerAuths: %w", err)
	}

	podmanAuths, err := getPodmanAuths()
	if err != nil {
		return nil, fmt.Errorf("getPodmanAuths: %w", err)
	}

	return append(dockerAuths, podmanAuths...), nil
}
