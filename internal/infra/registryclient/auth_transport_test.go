package registryclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

// newBearerAuthTestServer simulates a ghcr.io-like registry: every registry
// endpoint returns 401 with a Bearer challenge until a valid token is sent,
// and a /token endpoint issues tokens.
func newBearerAuthTestServer(t *testing.T, validToken string) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var tokenRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenRequests.Add(1)

			w.Header().Set("Content-Type", "application/json")

			token := validToken
			if token == "" {
				token = "anonymous-token"
			}

			fmt.Fprintf(w, `{"token": "%s", "expires_in": 300}`, token)
			return
		}

		if r.Header.Get("Authorization") == "Bearer "+validToken {
			w.Header().Set("Docker-Distribution-Api-Version", "registry/2.0")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"repositories": ["foo/bar"]}`)
			return
		}

		w.Header().Set(
			"WWW-Authenticate",
			`Bearer realm="`+r.Host+`/token",service="test-registry",scope="repository:foo/bar:pull"`,
		)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"errors":[{"code":"UNAUTHORIZED"}]}`)
	}))

	t.Cleanup(server.Close)

	return server, &tokenRequests
}

func TestAuthTransportBearerChallengeFlow(t *testing.T) {
	server, tokenRequests := newBearerAuthTestServer(t, "valid-token")

	at := newAuthTransport("octocat", "secret")
	at.SetBase(server.Client().Transport)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v2/foo/bar/tags/list", nil)
	require.NoError(t, err)

	resp, err := at.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, int32(1), tokenRequests.Load())

	// Second request for the same scope must reuse the cached token.
	req2, err := http.NewRequest(http.MethodGet, server.URL+"/v2/foo/bar/manifests/latest", nil)
	require.NoError(t, err)

	resp2, err := at.RoundTrip(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	require.Equal(t, http.StatusOK, resp2.StatusCode)
	require.Equal(t, int32(1), tokenRequests.Load())
}

func TestAuthTransportCredentialsSentToTokenEndpoint(t *testing.T) {
	var (
		tokenCalls   atomic.Int32
		basicAuthRaw atomic.Value
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenCalls.Add(1)
			basicAuthRaw.Store(r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"token": "the-token", "expires_in": 300}`)
			return
		}

		if r.Header.Get("Authorization") == "Bearer the-token" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "ok")
			return
		}

		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s/token",service="svc",scope="repository:a/b:pull"`, r.Host))
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	at := newAuthTransport("octocat", "secret")
	at.SetBase(server.Client().Transport)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v2/a/b/tags/list", nil)
	require.NoError(t, err)

	resp, err := at.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, int32(1), tokenCalls.Load())

	authHeader, _ := basicAuthRaw.Load().(string)
	require.Contains(t, authHeader, "Basic ")
}

func TestAuthTransportPassesThroughNonBearer401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	at := newAuthTransport("", "")
	at.SetBase(server.Client().Transport)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v2/", nil)
	require.NoError(t, err)

	resp, err := at.RoundTrip(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthTransportTokenEndpointFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"errors":[{"code":"DENIED"}]}`)
			return
		}

		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s/token",service="svc",scope="repository:a/b:pull"`, r.Host))
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	at := newAuthTransport("", "")
	at.SetBase(server.Client().Transport)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v2/a/b/tags/list", nil)
	require.NoError(t, err)

	_, err = at.RoundTrip(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fetch token")
}

func TestParseChallenge(t *testing.T) {
	header := `Bearer realm="https://ghcr.io/token",service="ghcr.io",scope="repository:user/image:pull"`

	params := parseChallenge(header)

	require.Equal(t, "https://ghcr.io/token", params["realm"])
	require.Equal(t, "ghcr.io", params["service"])
	require.Equal(t, "repository:user/image:pull", params["scope"])
}

func TestIsBearerChallenge(t *testing.T) {
	require.True(t, isBearerChallenge(`Bearer realm="https://example.com/token"`))
	require.True(t, isBearerChallenge(`bearer realm="https://example.com/token"`))
	require.False(t, isBearerChallenge(`Basic realm="example"`))
	require.False(t, isBearerChallenge(""))
}

func TestRequestScope(t *testing.T) {
	cases := []struct {
		rawURL string
		method string
		want   string
	}{
		{"http://example.com/v2/", http.MethodGet, ""},
		{"http://example.com/v2", http.MethodGet, ""},
		{"http://example.com/v2/foo/bar/tags/list", http.MethodGet, "repository:foo/bar:pull"},
		{"http://example.com/v2/foo/bar/manifests/latest", http.MethodGet, "repository:foo/bar:pull"},
		{"http://example.com/v2/foo/bar/blobs/sha256:abc", http.MethodGet, "repository:foo/bar:pull"},
		{"http://example.com/v2/foo/bar/manifests/sha256:abc", http.MethodDelete, "repository:foo/bar:pull,push,delete"},
	}

	for _, tc := range cases {
		u, err := url.Parse(tc.rawURL)
		require.NoError(t, err)

		req := &http.Request{Method: tc.method, URL: u}

		require.Equal(t, tc.want, requestScope(req), tc.rawURL)
	}
}

func TestRegistryClientRoundTripIntegration(t *testing.T) {
	server, _ := newBearerAuthTestServer(t, "valid-token")

	baseURL := strings.TrimPrefix(server.URL, "http://")

	client := New("http://"+baseURL, "octocat", "secret")

	repositories, err := client.ListRepositories(t.Context())
	require.NoError(t, err)
	require.Equal(t, []string{"foo/bar"}, repositories)

	require.NoError(t, client.Close())
}
