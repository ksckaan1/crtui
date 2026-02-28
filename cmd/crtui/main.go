package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrylist"
	"github.com/ksckaan1/crtui/config"
)

func main() {
	cfg, err := config.New(false)
	if err != nil {
		panic(err)
	}

	registryListScreen := registrylist.NewRegistryListScreenModel(cfg)

	program := tea.NewProgram(registryListScreen)
	_, err = program.Run()
	if err != nil {
		log.Fatalf("program.Run: %v", err)
	}
}
