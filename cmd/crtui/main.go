package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrylist"
	"github.com/ksckaan1/crtui/config"
	"github.com/ksckaan1/crtui/version"
)

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	v := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *showVersion || *v {
		fmt.Println(version.Version)
		os.Exit(0)
	}

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
