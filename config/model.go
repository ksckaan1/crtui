package config

type data struct {
	Auths map[string][]authData `yaml:"auths"`
}

type authData struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Auth struct {
	URL          string
	Username     string
	Password     string
	AutoDetected bool
}
