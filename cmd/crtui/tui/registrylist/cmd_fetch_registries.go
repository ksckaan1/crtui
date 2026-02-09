package registrylist

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksckaan1/crtui/internal/infra/registryclient"
)

type registryResult struct {
	index  int
	status string
}

func fetchRegistry(index int, registry *Registry) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		rc := registryclient.New(registry.URL, registry.Username, registry.Password)
		defer rc.Close()

		ri, err := rc.GetRegistryInfo(ctx)
		if err != nil {
			return registryResult{
				index:  index,
				status: "offline",
			}
		}

		if ri.AuthSuccess {
			return registryResult{
				index:  index,
				status: "online",
			}
		}

		return registryResult{
			index:  index,
			status: "unauth",
		}
	}
}
