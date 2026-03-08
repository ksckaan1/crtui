package updatecheckclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateCheckClient_GetLatestRelease(t *testing.T) {
	c := New()

	ctx := context.Background()

	release, err := c.GetLatestRelease(ctx, "ksckaan1", "crtui")
	require.NoError(t, err)

	t.Logf("%#v\n", release)
}
