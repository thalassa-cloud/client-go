package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func setupKubernetesTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/v1/kubernetes/clusters/cluster-123/upgradable-versions" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode([]KubernetesVersion{{Identity: "v1.32", Name: "1.32"}})
		case r.URL.Path == "/v1/kubernetes/versions/v1.31/upgradable-versions" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode([]KubernetesVersion{{Identity: "v1.32", Name: "1.32"}})
		case r.URL.Path == "/v1/kubernetes/clusters/cluster-123/container-images" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(ContainerImagesInUse{Images: []ContainerImage{{Image: "nginx", Tag: "1.25", Hash: "sha256:abc"}}})
		case r.URL.Path == "/v1/kubernetes/clusters/cluster-123/kubeconfig-sessions" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode([]KubernetesClusterSession{{Identity: "sess-1"}})
		case r.URL.Path == "/v1/kubernetes/clusters/cluster-123/kubeconfig-sessions/sess-1" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func newKubernetesTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	base, err := client.NewClient(client.WithBaseURL(serverURL), client.WithAuthCustom())
	require.NoError(t, err)
	k8s, err := New(base)
	require.NoError(t, err)
	return k8s
}

func TestGetUpgradableVersionsForCluster(t *testing.T) {
	server := setupKubernetesTestServer()
	defer server.Close()
	k8s := newKubernetesTestClient(t, server.URL)

	versions, err := k8s.GetUpgradableVersionsForCluster(context.Background(), "cluster-123")
	require.NoError(t, err)
	require.Len(t, versions, 1)
	assert.Equal(t, "v1.32", versions[0].Identity)

	_, err = k8s.GetUpgradableVersionsForCluster(context.Background(), "")
	assert.EqualError(t, err, "cluster identity is required")
}

func TestGetUpgradableVersions(t *testing.T) {
	server := setupKubernetesTestServer()
	defer server.Close()
	k8s := newKubernetesTestClient(t, server.URL)

	versions, err := k8s.GetUpgradableVersions(context.Background(), "v1.31")
	require.NoError(t, err)
	require.Len(t, versions, 1)

	_, err = k8s.GetUpgradableVersions(context.Background(), "")
	assert.EqualError(t, err, "version identity is required")
}

func TestListContainerImagesInUse(t *testing.T) {
	server := setupKubernetesTestServer()
	defer server.Close()
	k8s := newKubernetesTestClient(t, server.URL)

	images, err := k8s.ListContainerImagesInUse(context.Background(), "cluster-123")
	require.NoError(t, err)
	require.Len(t, images.Images, 1)
	assert.Equal(t, "nginx", images.Images[0].Image)
}

func TestKubeconfigSessions(t *testing.T) {
	server := setupKubernetesTestServer()
	defer server.Close()
	k8s := newKubernetesTestClient(t, server.URL)

	sessions, err := k8s.ListKubeconfigSessions(context.Background(), "cluster-123")
	require.NoError(t, err)
	require.Len(t, sessions, 1)
	assert.Equal(t, "sess-1", sessions[0].Identity)

	require.NoError(t, k8s.DeleteKubeconfigSession(context.Background(), "cluster-123", "sess-1"))
	assert.EqualError(t, k8s.DeleteKubeconfigSession(context.Background(), "cluster-123", ""), "session identity is required")
}
