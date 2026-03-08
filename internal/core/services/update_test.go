package services

import (
	"context"
	"testing"

	"github.com/ksckaan1/crtui/internal/infra/updatecheckclient"
	"github.com/stretchr/testify/require"
)

func TestUpdateService_CheckForUpdate(t *testing.T) {
	c := updatecheckclient.New()

	updateService := NewUpdateService(c, true, "dev")

	ctx := context.Background()

	release, err := updateService.CheckForUpdate(ctx)
	require.NoError(t, err)

	t.Logf("%#v\n", release)
}
