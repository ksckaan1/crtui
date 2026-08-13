package dockerhubclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListRepositories(t *testing.T) {
	var receivedAuthHeader string

	var server *httptest.Server

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/repositories/octocat/":
			w.Header().Set("Content-Type", "application/json")

			receivedAuthHeader = r.Header.Get("Authorization")

			if r.URL.Query().Get("page") == "2" {
				fmt.Fprint(w, `{"results":[{"name":"three","namespace":"octocat","repository_type":"image"}],"next":null}`)
				return
			}

			fmt.Fprint(w, `{"results":[{"name":"one","namespace":"octocat","repository_type":"image","is_private":false},{"name":"two","namespace":"octocat","repository_type":"image","is_private":true},{"name":"plugin-repo","namespace":"octocat","repository_type":"plugin"}],"next":"`+server.URL+`/v2/repositories/octocat/?page=2"}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "hunter2", server.URL)

	repos, err := client.ListRepositories(context.Background())
	require.NoError(t, err)
	require.Equal(t, []Repository{
		{Name: "one", Namespace: "octocat"},
		{Name: "two", Namespace: "octocat", IsPrivate: true},
		{Name: "three", Namespace: "octocat"},
	}, repos)
	require.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte("octocat:hunter2")), receivedAuthHeader)
}

func TestListRepositoriesAnonymous(t *testing.T) {
	var receivedAuthHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/repositories/octocat/":
			w.Header().Set("Content-Type", "application/json")
			receivedAuthHeader = r.Header.Get("Authorization")
			fmt.Fprint(w, `{"results":[{"name":"one","namespace":"octocat","repository_type":"image"}],"next":null}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "", server.URL)

	repos, err := client.ListRepositories(context.Background())
	require.NoError(t, err)
	require.Equal(t, []Repository{{Name: "one", Namespace: "octocat"}}, repos)
	require.Empty(t, receivedAuthHeader)
}

func TestListRepositoriesAnonymousPaginationLimit(t *testing.T) {
	var server *httptest.Server

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/repositories/octocat/":
			w.Header().Set("Content-Type", "application/json")

			if r.URL.Query().Get("page") == "2" {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, `{"message":"pagination offset too large for anonymous requests; sign in to page further"}`)
				return
			}

			fmt.Fprint(w, `{"results":[{"name":"one","namespace":"octocat","repository_type":"image"}],"next":"`+server.URL+`/v2/repositories/octocat/?page=2&page_size=25"}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "", server.URL)

	repos, err := client.ListRepositories(context.Background())
	require.NoError(t, err)
	require.Equal(t, []Repository{{Name: "one", Namespace: "octocat"}}, repos)
}

func TestListRepositoriesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/repositories/octocat/":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message":"namespace not found"}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "", server.URL)

	_, err := client.ListRepositories(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "docker hub api returned 404 Not Found")
}

func TestValidate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/repositories/octocat/":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"results":[],"next":null}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "hunter2", server.URL)

	err := client.Validate(context.Background())
	require.NoError(t, err)
}
