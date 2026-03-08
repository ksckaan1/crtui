package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/cmd/crtui/tui/registrylist"
	"github.com/ksckaan1/crtui/config"
	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/ksckaan1/crtui/internal/core/services"
	"github.com/ksckaan1/crtui/internal/infra/updatecheckclient"
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

	updateCheckClient := updatecheckclient.New()

	updateService := services.NewUpdateService(updateCheckClient, false, version.Version)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	release, _ := updateService.CheckForUpdate(ctx)

	var lastRelease *models.Release

	if release != nil && version.Version != release.Name {
		lastRelease = release
	}

	registryListScreen := registrylist.NewRegistryListScreenModel(cfg, lastRelease)

	program := tea.NewProgram(registryListScreen)
	_, err = program.Run()
	if err != nil {
		log.Fatalf("program.Run: %v", err)
	}
}
