package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/cli/cli/config"
	"github.com/stretchr/testify/require"
)

func TestIsDockerHubTokenKey(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"https://index.docker.io/v1/access-token", true},
		{"https://index.docker.io/v1/refresh-token", true},
		{"https://index.docker.io/v1/refresh-token/", true},
		{"https://registry-1.docker.io/access-token", true},
		{"access-token", false},
		{"refresh-token", false},
		{"https://index.docker.io/v1/", false},
		{"https://index.docker.io", false},
		{"ghcr.io", false},
		{"https://ghcr.io", false},
		{"registry.kaanksc.com", false},
		{"https://myregistry.example.com/repos", false},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, isDockerHubTokenKey(tc.key), tc.key)
	}
}

func TestListAuthsFromConfig(t *testing.T) {
	auth := func(user, pass string) string {
		return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	}

	configJSON := `{
  "auths": {
    "ghcr.io": {"auth": "` + auth("octocat", "ghp_token") + `"},
    "https://index.docker.io/v1/": {"auth": "` + auth("ksckaan1", "dckr_pat") + `"},
    "https://index.docker.io/v1/access-token": {"auth": "` + auth("ksckaan1", "jwt") + `"},
    "https://index.docker.io/v1/refresh-token": {"auth": "` + auth("ksckaan1", "refresh") + `"},
    "registry.kaanksc.com": {"auth": "` + auth("ksckaan1", "secret") + `"}
  }
}`

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	require.NoError(t, os.WriteFile(path, []byte(configJSON), 0o600))

	configFile, err := os.Open(path)
	require.NoError(t, err)
	defer configFile.Close()

	cfg, err := config.LoadFromReader(configFile)
	require.NoError(t, err)

	auths, err := listAuthsFromConfig(cfg)
	require.NoError(t, err)
	require.Len(t, auths, 3)

	byURL := map[string]*Auth{}
	for _, a := range auths {
		require.True(t, a.AutoDetected)
		byURL[a.URL] = a
	}

	require.Equal(t, "octocat", byURL["https://ghcr.io"].Username)
	require.Equal(t, "ghp_token", byURL["https://ghcr.io"].Password)

	require.Equal(t, "ksckaan1", byURL["https://registry-1.docker.io"].Username)
	require.Equal(t, "dckr_pat", byURL["https://registry-1.docker.io"].Password)

	require.Equal(t, "ksckaan1", byURL["https://registry.kaanksc.com"].Username)
	require.Equal(t, "secret", byURL["https://registry.kaanksc.com"].Password)
}

func TestConfig_New(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	err = cfg.Close()
	require.NoError(t, err)
}

func TestConfig_SetAndGetAuth(t *testing.T) {
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

func TestConfig_Update(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	oldURL := "https://old.com"
	oldUserName := "user1"
	err = cfg.SetAuth(oldURL, oldUserName, "pass1")
	require.NoError(t, err)

	newURL := "https://new.com"
	newUserName := "user2"
	err = cfg.UpdateAuth(oldURL, oldUserName, newURL, newUserName, "pass1")
	require.NoError(t, err)

	auth, err := cfg.GetAuth("https://new.com", "user2")
	require.NoError(t, err)
	require.Equal(t, "https://new.com", auth.URL)
	require.Equal(t, "user2", auth.Username)
	require.Equal(t, "pass1", auth.Password)
}

func TestConfig_SetAuthDuplicate(t *testing.T) {
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

func TestConfig_GetAuthNotFound(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	_, err = cfg.GetAuth("https://notfound.com", "user1")
	require.Error(t, err)
}

func TestConfig_GetAuthUsernameNotFound(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.SetAuth("https://test.com", "user1", "pass1")
	require.NoError(t, err)

	_, err = cfg.GetAuth("https://test.com", "nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "username not found")
}

func TestConfig_ListAuths(t *testing.T) {
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

func TestConfig_DeleteAuth(t *testing.T) {
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

func TestConfig_DeleteAuthNotFound(t *testing.T) {
	cfg, err := New(true, false)
	require.NoError(t, err)
	defer cfg.Close()

	err = cfg.DeleteAuth("https://notfound.com", "user1")
	require.NoError(t, err)
}

func TestConfig_EmptyUsername(t *testing.T) {
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

func TestConfig_ConfigFilePath(t *testing.T) {
	filename := "test_config.yaml"
	configFilePath, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), ".config", "crtui", filename))
	require.NoError(t, err)
	require.Contains(t, configFilePath, "crtui")
	require.Contains(t, configFilePath, "test_config.yaml")
}
