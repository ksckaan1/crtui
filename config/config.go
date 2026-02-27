package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
	"github.com/samber/lo"
	"github.com/zalando/go-keyring"
)

type Config struct {
	configFile        *os.File
	data              *data
	autoDetectedAuths []*Auth
}

func New(test bool) (*Config, error) {
	filename := "config.yaml"
	if test {
		filename = "test_config.yaml"
	}

	configFilePath, err := xdg.ConfigFile(filepath.Join("crtui", filename))
	if err != nil {
		return nil, fmt.Errorf("xdg.ConfigFile: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(configFilePath), 0755)
	if err != nil {
		return nil, fmt.Errorf("os.MkdirAll: %w", err)
	}

	configFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	b, err := io.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var d data

	if len(b) != 0 {
		err = yaml.Unmarshal(b, &d)
		if err != nil {
			return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
		}
	}

	autoDetectedAuths, err := listAutoDetectedAuths()
	if err != nil {
		return nil, fmt.Errorf("listAutoDetectedAuths: %w", err)
	}

	return &Config{
		configFile:        configFile,
		data:              &d,
		autoDetectedAuths: autoDetectedAuths,
	}, nil
}

func (c *Config) Close() error {
	return c.configFile.Close()
}

func (c *Config) ListAuths() ([]*Auth, error) {
	results := make([]*Auth, 0)

	for url, usernameMap := range c.data.Auths {
		for _, username := range usernameMap {
			password, err := keyring.Get("crtui", fmt.Sprintf("%s@%s", username, url))
			if err != nil {
				return nil, fmt.Errorf("keyring.Get: %w", err)
			}

			results = append(results, &Auth{
				URL:      url,
				Username: username,
				Password: password,
			})
		}
	}

	results = append(results, c.autoDetectedAuths...)

	slices.SortFunc(results, func(i, j *Auth) int {
		iTitle := fmt.Sprintf("%s-%s", i.URL, i.Username)
		jTitle := fmt.Sprintf("%s-%s", j.URL, j.Username)
		return strings.Compare(iTitle, jTitle)
	})

	return results, nil
}

func (c *Config) GetAuth(url, username string) (*Auth, error) {
	for _, auth := range c.autoDetectedAuths {
		if auth.URL == url && auth.Username == username {
			return auth, nil
		}
	}

	v, ok := c.data.Auths[url]
	if !ok {
		return nil, fmt.Errorf("url not found")
	}

	if !slices.Contains(v, username) {
		return nil, fmt.Errorf("username not found")
	}

	password, err := keyring.Get("crtui", fmt.Sprintf("%s@%s", username, url))
	if err != nil {
		return nil, fmt.Errorf("keyring.Get: %w", err)
	}

	return &Auth{
		URL:      url,
		Username: username,
		Password: password,
	}, nil
}

func (c *Config) SetAuth(url, username, password string) error {
	for _, auth := range c.autoDetectedAuths {
		if auth.URL == url && auth.Username == username {
			return fmt.Errorf("auth already exists")
		}
	}

	if c.data.Auths == nil {
		c.data.Auths = make(map[string][]string)
	}

	if _, ok := c.data.Auths[url]; !ok {
		c.data.Auths[url] = make([]string, 0)
	}

	if !slices.Contains(c.data.Auths[url], username) {
		c.data.Auths[url] = append(c.data.Auths[url], username)
	}

	err := keyring.Set("crtui", fmt.Sprintf("%s@%s", username, url), password)
	if err != nil {
		return fmt.Errorf("keyring.Set: %w", err)
	}

	err = c.save()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (c *Config) DeleteAuth(url, username string) error {
	if c.data.Auths == nil {
		return nil
	}

	if _, ok := c.data.Auths[url]; !ok {
		return nil
	}

	c.data.Auths[url] = lo.Filter(c.data.Auths[url], func(item string, _ int) bool {
		return item != username
	})

	err := keyring.Delete("crtui", fmt.Sprintf("%s@%s", username, url))
	if err != nil {
		return fmt.Errorf("keyring.Delete: %w", err)
	}

	err = c.save()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (c *Config) save() error {
	err := c.configFile.Truncate(0)
	if err != nil {
		return fmt.Errorf("configFile.Truncate: %w", err)
	}

	_, err = c.configFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("configFile.Seek: %w", err)
	}

	err = yaml.NewEncoder(c.configFile).Encode(c.data)
	if err != nil {
		return fmt.Errorf("yaml.NewEncoder: %w", err)
	}

	return nil
}
