package registryclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type tokenResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type cachedToken struct {
	token     string
	expiresAt time.Time
}

// authTransport is an http.RoundTripper that transparently implements the OCI
// Distribution token authentication flow. When the registry responds with a
// 401 and a `WWW-Authenticate: Bearer` challenge, it exchanges the configured
// credentials for a token at the challenge realm and retries the request once
// with that token. Tokens are cached per scope so that subsequent requests do
// not need another round trip.
type authTransport struct {
	base     http.RoundTripper
	username string
	password string

	mu     sync.Mutex
	tokens map[string]cachedToken
}

func newAuthTransport(username, password string) *authTransport {
	return &authTransport{
		username: username,
		password: password,
		tokens:   make(map[string]cachedToken),
	}
}

func (t *authTransport) SetBase(base http.RoundTripper) {
	t.base = base
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.base == nil {
		return nil, fmt.Errorf("registryclient: auth transport base is not set")
	}

	key := requestScope(req)

	if token, ok := t.getToken(key); ok {
		cloned := cloneWithAuthorization(req, token)

		resp, err := t.base.RoundTrip(cloned)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}

		if challenge := resp.Header.Get("WWW-Authenticate"); isBearerChallenge(challenge) {
			return t.retryWithChallenge(req, resp, challenge)
		}

		return resp, nil
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	challenge := resp.Header.Get("WWW-Authenticate")
	if !isBearerChallenge(challenge) {
		return resp, nil
	}

	return t.retryWithChallenge(req, resp, challenge)
}

func (t *authTransport) retryWithChallenge(req *http.Request, unauthorized *http.Response, challenge string) (*http.Response, error) {
	params := parseChallenge(challenge)

	realm := params["realm"]
	if realm == "" {
		return unauthorized, nil
	}

	if !strings.Contains(realm, "://") && req.URL != nil {
		realm = req.URL.Scheme + "://" + realm
	}

	scope := params["scope"]

	token, expiresAt, err := t.fetchToken(req.Context(), realm, params["service"], scope)
	if err != nil {
		unauthorized.Body.Close()
		return nil, fmt.Errorf("registryclient: fetch token: %w", err)
	}

	t.setToken(scope, token, expiresAt)

	unauthorized.Body.Close()

	return t.base.RoundTrip(cloneWithAuthorization(req, token))
}

func (t *authTransport) fetchToken(ctx context.Context, realm, service, scope string) (string, time.Time, error) {
	tokenURL, err := url.Parse(realm)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("parse realm %q: %w", realm, err)
	}

	query := tokenURL.Query()
	if service != "" {
		query.Set("service", service)
	}
	if scope != "" {
		query.Set("scope", scope)
	}
	tokenURL.RawQuery = query.Encode()

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL.String(), nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("new token request: %w", err)
	}
	tokenReq.Header.Set("Accept", "application/json")

	if t.username != "" || t.password != "" {
		tokenReq.SetBasicAuth(t.username, t.password)
	}

	tokenResp, err := t.base.RoundTrip(tokenReq)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("token endpoint: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(tokenResp.Body, 512))
		return "", time.Time{}, fmt.Errorf(
			"token endpoint returned %s: %s",
			tokenResp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var parsed tokenResponse

	if err := json.NewDecoder(tokenResp.Body).Decode(&parsed); err != nil {
		return "", time.Time{}, fmt.Errorf("decode token response: %w", err)
	}

	token := parsed.Token
	if token == "" {
		token = parsed.AccessToken
	}
	if token == "" {
		return "", time.Time{}, fmt.Errorf("token endpoint returned an empty token")
	}

	expiresAt := time.Time{}
	if parsed.ExpiresIn > 0 {
		expiresAt = time.Now().Add(time.Duration(parsed.ExpiresIn) * time.Second)
	}

	return token, expiresAt, nil
}

func (t *authTransport) getToken(scope string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	entry, ok := t.tokens[scope]
	if !ok {
		return "", false
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(t.tokens, scope)
		return "", false
	}

	return entry.token, true
}

func (t *authTransport) setToken(scope, token string, expiresAt time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tokens[scope] = cachedToken{
		token:     token,
		expiresAt: expiresAt,
	}
}

func cloneWithAuthorization(req *http.Request, token string) *http.Request {
	cloned := req.Clone(req.Context())
	cloned.Header = req.Header.Clone()
	cloned.Header.Set("Authorization", "Bearer "+token)
	return cloned
}

func isBearerChallenge(header string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(header)), "bearer ")
}

func parseChallenge(header string) map[string]string {
	params := make(map[string]string)

	_, rest, ok := strings.Cut(strings.TrimSpace(header), " ")
	if !ok {
		return params
	}

	for _, part := range strings.Split(rest, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}

		params[strings.ToLower(strings.TrimSpace(key))] = strings.Trim(strings.TrimSpace(value), `"`)
	}

	return params
}

// requestScope builds a best-effort token cache key from the request path and
// method. It is only used to reuse a previously issued token; the authoritative
// scope always comes from the registry challenge, so a wrong guess only causes
// a fresh token exchange.
func requestScope(req *http.Request) string {
	path := strings.TrimPrefix(req.URL.Path, "/v2/")
	path = strings.TrimPrefix(path, "/")

	if path == "" || path == "v2" {
		return ""
	}

	repoParts := make([]string, 0)
	for _, part := range strings.Split(path, "/") {
		if part == "manifests" || part == "blobs" || part == "tags" {
			break
		}
		repoParts = append(repoParts, part)
	}

	repo := strings.Join(repoParts, "/")
	if repo == "" {
		return ""
	}

	action := "pull"
	if req.Method == http.MethodDelete {
		action = "pull,push,delete"
	}

	return fmt.Sprintf("repository:%s:%s", repo, action)
}
