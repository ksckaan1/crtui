package config

type data struct {
	Auths map[string][]string `yaml:"auths"`
}

type Auth struct {
	URL          string
	Username     string
	Password     string
	AutoDetected bool
}
