package config

import (
	"encoding/base64"
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
	isKeyringEnabled  bool
}

func New(test bool, keyringEnabled bool) (*Config, error) {
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
			return nil, fmt.Errorf("yaml.Unmarshal: %w (%s)", err, configFilePath)
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
		isKeyringEnabled:  keyringEnabled,
	}, nil
}

func (c *Config) Close() error {
	return c.configFile.Close()
}

func (c *Config) ListAuths() ([]*Auth, error) {
	results := make([]*Auth, 0)

	for url, usernameMap := range c.data.Auths {
		for _, authData := range usernameMap {
			var (
				password string
				err      error
			)

			if authData.Username == "" {
				results = append(results, &Auth{
					URL:      url,
					Username: "",
					Password: "",
				})
				continue
			}

			if authData.Password == "" && c.isKeyringEnabled {
				password, err = keyring.Get("crtui", fmt.Sprintf("%s@%s", authData.Username, url))
				if err != nil {
					return nil, fmt.Errorf("keyring.Get: %w", err)
				}
			} else {
				passwordBytes, err := base64.StdEncoding.DecodeString(authData.Password)
				if err != nil {
					return nil, fmt.Errorf("base64.StdEncoding.DecodeString: %w", err)
				}
				password = string(passwordBytes)
			}

			results = append(results, &Auth{
				URL:      url,
				Username: authData.Username,
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

	foundAuth, ok := lo.Find(v, func(item authData) bool {
		return item.Username == username

	})
	if !ok {
		return nil, fmt.Errorf("username not found")
	}

	if foundAuth.Username == "" {
		return &Auth{
			URL: url,
		}, nil
	}

	var (
		password string
		err      error
	)

	if foundAuth.Password == "" && c.isKeyringEnabled {
		password, err = keyring.Get("crtui", fmt.Sprintf("%s@%s", username, url))
		if err != nil {
			return nil, fmt.Errorf("keyring.Get: %w", err)
		}
	} else {
		passwordBytes, err := base64.StdEncoding.DecodeString(foundAuth.Password)
		if err != nil {
			return nil, fmt.Errorf("base64.StdEncoding.DecodeString: %w", err)
		}
		password = string(passwordBytes)
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
		c.data.Auths = make(map[string][]authData)
	}

	if _, ok := c.data.Auths[url]; !ok {
		c.data.Auths[url] = make([]authData, 0)
	}

	err := keyring.Set("crtui", fmt.Sprintf("%s@%s", username, url), password)
	if err == nil && c.isKeyringEnabled {
		password = ""
	} else {
		password = base64.StdEncoding.EncodeToString([]byte(password))
	}

	_, index, ok := lo.FindIndexOf(c.data.Auths[url], func(item authData) bool {
		return item.Username == username
	})

	if !ok {
		c.data.Auths[url] = append(c.data.Auths[url], authData{Username: username, Password: password})
	} else {
		c.data.Auths[url][index].Password = password
	}

	err = c.save()
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (c *Config) UpdateAuth(oldURL, oldUsername, newURL, newUserName, password string) error {
	err := c.DeleteAuth(oldURL, oldUsername)
	if err != nil {
		return fmt.Errorf("delete auth: %w", err)
	}

	err = c.SetAuth(newURL, newUserName, password)
	if err != nil {
		return fmt.Errorf("set auth: %w", err)
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

	if _, ok := lo.Find(c.data.Auths[url], func(item authData) bool {
		return item.Username == username
	}); !ok {
		return nil
	}

	c.data.Auths[url] = lo.Filter(c.data.Auths[url], func(item authData, _ int) bool {
		return item.Username != username
	})

	if c.isKeyringEnabled {
		keyring.Delete("crtui", fmt.Sprintf("%s@%s", username, url))
	}

	for group, elems := range c.data.Auths {
		if len(elems) == 0 {
			delete(c.data.Auths, group)
		}
	}

	err := c.save()
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
