package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/ksckaan1/crtui/internal/core/enums/registrytype"
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
		dockerAuths = append(dockerAuths, newAutoDetectedAuth(key, authConfig.Username, authConfig.Password))
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
		dockerAuths = append(dockerAuths, newAutoDetectedAuth(key, authConfig.Username, authConfig.Password))
	}

	return dockerAuths, nil
}

func newAutoDetectedAuth(key, username, password string) *Auth {
	if !strings.HasPrefix(key, "http://") && !strings.HasPrefix(key, "https://") {
		key = "https://" + key
	}

	return &Auth{
		URL:          registrytype.Normalize(key),
		Username:     username,
		Password:     password,
		AutoDetected: true,
	}
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
