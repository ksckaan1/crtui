package models

import (
	"time"

	"github.com/ksckaan1/crtui/internal/core/enums/registrystatus"
)

type Registry struct {
	BaseURL       string
	SupportsHTTP3 bool
	Username      string
	Status        registrystatus.RegistryStatus
	Took          time.Duration
}

type TagList struct {
	Name string
	Tags []string
}

type Tag struct {
	RepositoryName string
	TagName        string
	Platforms      []Platform
}

type Platform struct {
	Config      Config
	Layers      []Layer
	Annotations map[string]string
}

type Config struct {
	Digest       string
	Architecture string
	OS           string
	Author       string
	Created      time.Time
	History      []History
	Entrypoint   []string
	Env          []string
	Labels       map[string]string
	User         string
}

type History struct {
	Author    string
	Created   time.Time
	CreatedBy string
	Comment   string
}

type Layer struct {
	MediaType string
	Size      int
	Digest    string
}
