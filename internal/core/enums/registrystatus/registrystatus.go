package registrystatus

type RegistryStatus string

const (
	Loading RegistryStatus = "loading"
	Invalid RegistryStatus = "invalid"
	Unauth  RegistryStatus = "unauth"
	Offline RegistryStatus = "offline"
	Online  RegistryStatus = "online"
)
