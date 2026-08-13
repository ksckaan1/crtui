package githubclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListContainerPackagesUser(t *testing.T) {
	var authorizationHeader string

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader = r.Header.Get("Authorization")

		switch r.URL.Path {
		case "/orgs/octocat":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message":"Not Found"}`)
			return
		case "/users/octocat/packages":
			w.Header().Set("Content-Type", "application/json")

			if r.URL.Query().Get("page") == "2" {
				fmt.Fprint(w, `[{"name":"c","package_type":"container"}]`)
				return
			}

			w.Header().Set("Link", fmt.Sprintf(
				`<%s/users/octocat/packages?page=2&package_type=container&per_page=100>; rel="next", <%s/users/octocat/packages?page=2&package_type=container&per_page=100>; rel="last"`,
				server.URL,
				server.URL,
			))

			fmt.Fprint(w, `[{"name":"a","package_type":"container"},{"name":"b","package_type":"container"},{"name":"x","package_type":"npm"}]`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "ghp_123", server.URL)

	repos, err := client.ListContainerPackages(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"octocat/a", "octocat/b", "octocat/c"}, repos)
	require.Equal(t, "Bearer ghp_123", authorizationHeader)
}

func TestListContainerPackagesOrg(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orgs/kubernetes":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"login":"kubernetes","type":"Organization"}`)
			return
		case "/orgs/kubernetes/packages":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[{"name":"k8s","package_type":"container"}]`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("kubernetes", "ghp_123", server.URL)

	repos, err := client.ListContainerPackages(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"kubernetes/k8s"}, repos)
}

func TestListContainerPackagesRequiresAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orgs/octocat":
			w.WriteHeader(http.StatusNotFound)
			return
		case "/users/octocat/packages":
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"message":"Requires authentication"}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "", server.URL)

	_, err := client.ListContainerPackages(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "401")
}

func TestValidate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/orgs/octocat":
			w.WriteHeader(http.StatusNotFound)
			return
		case "/users/octocat/packages":
			if r.URL.Query().Get("per_page") != "1" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[]`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := NewWithBaseURL("octocat", "ghp_123", server.URL)

	require.NoError(t, client.Validate(context.Background()))
}

func TestNextLink(t *testing.T) {
	link := `<https://api.github.com/users/octocat/packages?page=2>; rel="next", <https://api.github.com/users/octocat/packages?page=5>; rel="last"`

	require.Equal(t, "https://api.github.com/users/octocat/packages?page=2", nextLink(link))
	require.Equal(t, "", nextLink(""))
	require.Equal(t, "", nextLink(`<https://api.github.com/a>; rel="last"`))
}
