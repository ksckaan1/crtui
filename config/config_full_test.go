package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	err = cfg.Close()
	require.NoError(t, err)
}

func TestSetAndGetAuth(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.SetAuth("https://test.com", "user1", "pass1")
	require.NoError(t, err)

	auth, err := cfg.GetAuth("https://test.com", "user1")
	require.NoError(t, err)
	require.Equal(t, "https://test.com", auth.URL)
	require.Equal(t, "user1", auth.Username)
	require.Equal(t, "pass1", auth.Password)
}

func TestSetAuthDuplicate(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.SetAuth("https://test.com", "user1", "pass1")
	require.NoError(t, err)

	err = cfg.SetAuth("https://test.com", "user1", "pass2")
	require.NoError(t, err)

	auth, err := cfg.GetAuth("https://test.com", "user1")
	require.NoError(t, err)
	require.Equal(t, "pass2", auth.Password)
}

func TestGetAuthNotFound(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	_, err = cfg.GetAuth("https://notfound.com", "user1")
	require.Error(t, err)
}

func TestGetAuthUsernameNotFound(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.SetAuth("https://test.com", "user1", "pass1")
	require.NoError(t, err)

	_, err = cfg.GetAuth("https://test.com", "nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "username not found")
}

func TestListAuths(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.SetAuth("https://test1.com", "user1", "pass1")
	require.NoError(t, err)

	err = cfg.SetAuth("https://test2.com", "user2", "pass2")
	require.NoError(t, err)

	auths, err := cfg.ListAuths()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(auths), 2)
}

func TestDeleteAuth(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.SetAuth("https://test.com", "user1", "pass1")
	require.NoError(t, err)

	err = cfg.DeleteAuth("https://test.com", "user1")
	require.NoError(t, err)

	_, err = cfg.GetAuth("https://test.com", "user1")
	require.Error(t, err)
}

func TestDeleteAuthNotFound(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.DeleteAuth("https://notfound.com", "user1")
	require.NoError(t, err)
}

func TestAuthPasswordEncoding(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	originalPassword := "testPassword123!@#"
	err = cfg.SetAuth("https://test.com", "user1", originalPassword)
	require.NoError(t, err)

	auth, err := cfg.GetAuth("https://test.com", "user1")
	require.NoError(t, err)
	require.Equal(t, originalPassword, auth.Password)
}

func TestEmptyUsername(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	auths, err := cfg.ListAuths()
	require.NoError(t, err)

	for _, auth := range auths {
		if auth.Username == "" {
			require.Equal(t, "", auth.Password)
		}
	}
}

func TestConfigFilePath(t *testing.T) {
	filename := "test_config.yaml"
	configFilePath, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), ".config", "crtui", filename))
	require.NoError(t, err)
	require.Contains(t, configFilePath, "crtui")
	require.Contains(t, configFilePath, "test_config.yaml")
}

func TestDataStructure(t *testing.T) {
	d := &data{
		Auths: map[string][]authData{
			"https://example.com": {
				{Username: "user1", Password: base64.StdEncoding.EncodeToString([]byte("pass1"))},
			},
		},
	}
	require.NotNil(t, d.Auths)
	require.Contains(t, d.Auths, "https://example.com")
	require.Len(t, d.Auths["https://example.com"], 1)
}
