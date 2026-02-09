package registryclient

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryClient(t *testing.T) {
	client := New("https://cr.kaanksc.com", "ksckaan1", "@thXbdjxQU1IcMVv9d1y")
	defer client.Close()

	ctx := context.Background()
	_, err := client.GetRegistryInfo(ctx)
	require.NoError(t, err)

	// t.Logf("%#v", ri)

	// repositories, err := client.ListRepositories(ctx)
	// require.NoError(t, err)

	// t.Log(repositories)

	// tagList, err := client.ListTags(ctx, "piriapi/api")
	// require.NoError(t, err)

	// t.Logf("%#v", tagList)

	tag, err := client.GetTag(ctx, "piriapi/api", "0.0.4")
	require.NoError(t, err)

	jsonBytes, err := json.MarshalIndent(tag, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal tag: %v", err)
	}
	t.Logf("%s", jsonBytes)

	// err = client.DeleteTag(ctx, "piriapi/api", "0.0.3")
	// require.NoError(t, err)
}
