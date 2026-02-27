package registrylist

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
)

type registryResult struct {
	index         int
	status        registrystatus.RegistryStatus
	supportsHTTP3 bool
}

func fetchRegistry(index int, registry *Registry) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		rc := registryclient.New(registry.URL, registry.Username, registry.Password)
		defer rc.Close()

		ri, err := rc.GetRegistryInfo(ctx)
		if err != nil {
			return registryResult{
				index:         index,
				status:        registrystatus.Offline,
				supportsHTTP3: false,
			}
		}

		return registryResult{
			index:         index,
			status:        ri.Status,
			supportsHTTP3: ri.SupportsHTTP3,
		}
	}
}
