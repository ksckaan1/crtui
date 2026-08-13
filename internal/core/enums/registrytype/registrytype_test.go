package registrytype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromURL(t *testing.T) {
	cases := []struct {
		rawURL string
		want   RegistryType
	}{
		{"https://ghcr.io", GitHub},
		{"https://ghcr.io/", GitHub},
		{"http://ghcr.io", GitHub},
		{"GHCR.IO", GitHub},
		{"https://registry.example.com", Docker},
		{"http://localhost:5000", Docker},
		{"not a url", Docker},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, FromURL(tc.rawURL), tc.rawURL)
	}
}
