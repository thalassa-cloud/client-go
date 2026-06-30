package projects

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	c, err := client.NewClient(
		client.WithBaseURL(serverURL),
		client.WithAuthCustom(),
		client.WithOrganisation("acme"),
		client.WithProject("prj-should-not-be-sent"),
	)
	require.NoError(t, err)
	projectsClient, err := New(c)
	require.NoError(t, err)
	return projectsClient
}

func assertOrgScopedHeaders(t *testing.T, r *http.Request) {
	t.Helper()
	assert.Equal(t, "acme", r.Header.Get("X-Organisation-Identity"))
	assert.Empty(t, r.Header.Get("X-Project-Identity"))
}

func TestListProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects", r.URL.Path)
		assert.Equal(t, "production", r.URL.Query().Get("slug"))
		assert.Equal(t, "prod", r.URL.Query().Get("matchLabels[env]"))
		assertOrgScopedHeaders(t, r)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]Project{
			{
				Identity:    "prj-abc123",
				Name:        "testproject",
				Slug:        "testproject",
				Description: "Root project",
			},
		}))
	}))
	defer server.Close()

	projectsClient := newTestClient(t, server.URL)

	projects, err := projectsClient.ListProjects(context.Background(), &ListProjectsRequest{
		Filters: []filters.Filter{
			&filters.FilterKeyValue{Key: filters.FilterSlug, Value: "production"},
			&filters.LabelFilter{MatchLabels: map[string]string{"env": "prod"}},
		},
	})
	require.NoError(t, err)
	require.Len(t, projects, 1)
	assert.Equal(t, "prj-abc123", projects[0].Identity)
	assert.Equal(t, "testproject", projects[0].Slug)
}

func TestGetProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects/production", r.URL.Path)
		assertOrgScopedHeaders(t, r)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(Project{
			Identity: "prj-abc123",
			Name:     "Production",
			Slug:     "production",
			ParentProject: &ProjectRef{
				Identity: "prj-testproject",
				Name:     "testproject",
				Slug:     "testproject",
			},
		}))
	}))
	defer server.Close()

	projectsClient := newTestClient(t, server.URL)

	project, err := projectsClient.GetProject(context.Background(), "production")
	require.NoError(t, err)
	assert.Equal(t, "prj-abc123", project.Identity)
	require.NotNil(t, project.ParentProject)
	assert.Equal(t, "prj-testproject", project.ParentProject.Identity)
}

func TestCreateProject(t *testing.T) {
	parent := "prj-testproject"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects", r.URL.Path)
		assertOrgScopedHeaders(t, r)

		var body CreateProjectRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Team A", body.Name)
		require.NotNil(t, body.ParentProjectIdentity)
		assert.Equal(t, parent, *body.ParentProjectIdentity)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(Project{
			Identity:      "prj-abc123",
			Name:          body.Name,
			Slug:          "abc123",
			ObjectVersion: 1,
			CreatedAt:     time.Now().UTC(),
		}))
	}))
	defer server.Close()

	projectsClient := newTestClient(t, server.URL)

	project, err := projectsClient.CreateProject(context.Background(), CreateProjectRequest{
		Name:                  "Team A",
		Labels:                map[string]string{},
		Annotations:           map[string]string{},
		ParentProjectIdentity: &parent,
	})
	require.NoError(t, err)
	assert.Equal(t, "prj-abc123", project.Identity)
	assert.Equal(t, "abc123", project.Slug)
}

func TestUpdateProjectClearParent(t *testing.T) {
	emptyParent := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v1/projects/prj-abc123", r.URL.Path)
		assertOrgScopedHeaders(t, r)

		var body UpdateProjectRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Team A", body.Name)
		require.NotNil(t, body.ParentProjectIdentity)
		assert.Equal(t, "", *body.ParentProjectIdentity)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(Project{
			Identity: "prj-abc123",
			Name:     body.Name,
			Slug:     "abc123",
		}))
	}))
	defer server.Close()

	projectsClient := newTestClient(t, server.URL)

	project, err := projectsClient.UpdateProject(context.Background(), "prj-abc123", UpdateProjectRequest{
		Name:                  "Team A",
		Description:           "Renamed",
		Labels:                map[string]string{},
		Annotations:           map[string]string{},
		ParentProjectIdentity: &emptyParent,
	})
	require.NoError(t, err)
	assert.Equal(t, "abc123", project.Slug)
}

func TestDeleteProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/projects/prj-abc123", r.URL.Path)
		assertOrgScopedHeaders(t, r)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	projectsClient := newTestClient(t, server.URL)

	err := projectsClient.DeleteProject(context.Background(), "prj-abc123")
	require.NoError(t, err)
}

func TestGetProjectNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]string{
			"message": "project not found",
		}))
	}))
	defer server.Close()

	projectsClient := newTestClient(t, server.URL)

	_, err := projectsClient.GetProject(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, client.IsNotFound(err))
	assert.Contains(t, err.Error(), "project not found")
}

func TestCreateProjectValidation(t *testing.T) {
	projectsClient := newTestClient(t, "http://example.com")

	_, err := projectsClient.CreateProject(context.Background(), CreateProjectRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}
