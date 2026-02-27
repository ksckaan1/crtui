package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	cfg, err := New(true)
	require.NoError(t, err)

	err = cfg.SetAuth("https://test.com", "ksckaan1", "abc123")
	require.NoError(t, err)

	err = cfg.SetAuth("https://test.com", "old", "abc123")
	require.NoError(t, err)

	auth, err := cfg.GetAuth("https://test.com", "ksckaan1")
	require.NoError(t, err)

	t.Logf("%#v\n", auth)

	auths, err := cfg.ListAuths()
	require.NoError(t, err)

	for i, auth := range auths {
		t.Logf("%d - %#v\n", i, auth)
	}

	err = cfg.DeleteAuth("https://test.com", "old")
	require.NoError(t, err)

	auths, err = cfg.ListAuths()
	require.NoError(t, err)

	for i, auth := range auths {
		t.Logf("%d - %#v\n", i, auth)
	}
}
