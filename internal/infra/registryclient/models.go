package registryclient

import "time"

type imageIndexV1 struct {
	SchemaVersion int                    `json:"schemaVersion"`
	MediaType     string                 `json:"mediaType"`
	Manifests     []imageIndexV1Manifest `json:"manifests"`
	Annotations   map[string]string      `json:"annotations"`
}

type imageIndexV1Manifest struct {
	MediaType string                       `json:"mediaType"`
	Size      int                          `json:"size"`
	Digest    string                       `json:"digest"`
	Platform  imageIndexV1ManifestPlatform `json:"platform"`
}

type imageIndexV1ManifestPlatform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

type imageManifestV1 struct {
	SchemaVersion int                    `json:"schemaVersion"`
	MediaType     string                 `json:"mediaType"`
	Config        imageManifestV1Config  `json:"config"`
	Layers        []imageManifestV1Layer `json:"layers"`
	Annotations   map[string]string      `json:"annotations"`
}

type imageManifestV1Config struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type imageManifestV1Layer struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type blob struct {
	Architecture string        `json:"architecture"`
	OS           string        `json:"os"`
	Author       string        `json:"author"`
	Created      time.Time     `json:"created"`
	History      []blobHistory `json:"history"`
	RootFS       blobRootFS    `json:"rootfs"`
	Config       blobConfig    `json:"config"`
}

type blobHistory struct {
	Author     string    `json:"author"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Comment    string    `json:"comment"`
	EmptyLayer bool      `json:"empty_layer"`
}

type blobRootFS struct {
	Type    string   `json:"type"`
	DiffIDs []string `json:"diff_ids"`
}

type blobConfig struct {
	Entrypoint []string          `json:"Entrypoint"`
	Cmd        []string          `json:"Cmd"`
	Env        []string          `json:"Env"`
	WorkingDir string            `json:"WorkingDir"`
	Labels     map[string]string `json:"Labels"`
	User       string            `json:"User"`
}
