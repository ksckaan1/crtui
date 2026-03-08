package registryclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	client := New("https://registry.example.com", "user", "pass")
	require.NotNil(t, client)
	require.Equal(t, "https://registry.example.com", client.BaseURL)
	require.Equal(t, "user", client.username)
	require.Equal(t, "pass", client.password)
	err := client.Close()
	require.NoError(t, err)
}

func TestNewWithoutCredentials(t *testing.T) {
	client := New("https://registry.example.com", "", "")
	require.NotNil(t, client)
	require.Equal(t, "https://registry.example.com", client.BaseURL)
	err := client.Close()
	require.NoError(t, err)
}

func TestSetTransportHTTP3(t *testing.T) {
	client := New("https://registry.example.com", "user", "pass")
	require.NotNil(t, client)
	require.False(t, client.supportsHTTP3)

	client.SetTransport(true)
	require.True(t, client.supportsHTTP3)

	client.SetTransport(false)
	require.False(t, client.supportsHTTP3)

	err := client.Close()
	require.NoError(t, err)
}

func TestSetTransportNoChange(t *testing.T) {
	client := New("https://registry.example.com", "user", "pass")
	require.NotNil(t, client)

	client.SetTransport(true)
	require.True(t, client.supportsHTTP3)

	client.SetTransport(true)
	require.True(t, client.supportsHTTP3)

	err := client.Close()
	require.NoError(t, err)
}
