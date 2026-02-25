package main

import (
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/docker/cli/cli/config"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrylist"
)

func main() {
	cfg, err := config.Load(config.Dir())
	if err != nil {
		log.Fatalf("config.Load: %v", err)
	}

	creds, err := cfg.GetAllCredentials()
	if err != nil {
		log.Fatalf("cfg.GetAllCredentials: %v", err)
	}

	registries := make([]*registrylist.Registry, 0)

	for key, authConfig := range creds {
		if strings.HasPrefix(key, "https://index.docker.io/v1/") {
			continue
		}

		registries = append(registries, &registrylist.Registry{
			URL:      key,
			Username: authConfig.Username,
			Password: authConfig.Password,
		})
	}

	registryListModel := registrylist.NewModel(registries)

	program := tea.NewProgram(registryListModel)
	_, err = program.Run()
	if err != nil {
		log.Fatalf("program.Run: %v", err)
	}
}
