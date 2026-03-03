package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrylist"
	"github.com/ksckaan1/crtui/config"
)

func main() {
	cfg, err := config.New(false, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	registryListScreen := registrylist.NewRegistryListScreenModel(cfg)

	program := tea.NewProgram(registryListScreen)
	_, err = program.Run()
	if err != nil {
		log.Fatalf("program.Run: %v", err)
	}
}
