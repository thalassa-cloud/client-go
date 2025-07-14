package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestListKubernetesClusterRoles(t *testing.T) {
	tests := []struct {
		name           string
		request        *ListKubernetesClusterRolesRequest
		serverResponse []KubernetesClusterRole
		serverStatus   int
		expectError    bool
		expectedCount  int
	}{
		{
			name:           "successful list without filters",
			request:        nil,
			serverResponse: []KubernetesClusterRole{},
			serverStatus:   http.StatusOK,
			expectError:    false,
			expectedCount:  0,
		},
		{
			name: "successful list with filters",
			request: &ListKubernetesClusterRolesRequest{
				Filters: []filters.Filter{
					&filters.FilterKeyValue{
						Key:   "name",
						Value: "test-role",
					},
				},
			},
			serverResponse: []KubernetesClusterRole{
				{
					Identity:    "role-1",
					Name:        "test-role",
					Slug:        "test-role",
					Description: "Test role",
					CreatedAt:   time.Now(),
				},
			},
			serverStatus:  http.StatusOK,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:           "server error",
			request:        nil,
			serverResponse: []KubernetesClusterRole{},
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" && strings.Contains(r.URL.Path, KubernetesClusterRoleEndpoint) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					if tt.serverStatus == http.StatusOK {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			result, err := c.ListKubernetesClusterRoles(context.Background(), tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}
}

func TestCreateKubernetesClusterRole(t *testing.T) {
	tests := []struct {
		name           string
		createRequest  CreateKubernetesClusterRoleRequest
		serverResponse *KubernetesClusterRole
		serverStatus   int
		expectError    bool
	}{
		{
			name: "successful creation",
			createRequest: CreateKubernetesClusterRoleRequest{
				Name:        "test-role",
				Description: "Test role description",
				Labels: map[string]string{
					"environment": "test",
				},
				Annotations: map[string]string{
					"created-by": "test",
				},
			},
			serverResponse: &KubernetesClusterRole{
				Identity:    "role-1",
				Name:        "test-role",
				Slug:        "test-role",
				Description: "Test role description",
				CreatedAt:   time.Now(),
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name: "server error",
			createRequest: CreateKubernetesClusterRoleRequest{
				Name:        "test-role",
				Description: "Test role description",
			},
			serverResponse: nil,
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" && strings.Contains(r.URL.Path, KubernetesClusterRoleEndpoint) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					if tt.serverStatus == http.StatusOK && tt.serverResponse != nil {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			result, err := c.CreateKubernetesClusterRole(context.Background(), tt.createRequest)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestGetKubernetesClusterRole(t *testing.T) {
	tests := []struct {
		name           string
		identity       string
		serverResponse *KubernetesClusterRole
		serverStatus   int
		expectError    bool
	}{
		{
			name:     "successful retrieval",
			identity: "role-1",
			serverResponse: &KubernetesClusterRole{
				Identity:    "role-1",
				Name:        "test-role",
				Slug:        "test-role",
				Description: "Test role description",
				CreatedAt:   time.Now(),
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:           "not found",
			identity:       "non-existent",
			serverResponse: nil,
			serverStatus:   http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" && strings.HasPrefix(r.URL.Path, KubernetesClusterRoleEndpoint+"/") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					if tt.serverStatus == http.StatusOK && tt.serverResponse != nil {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			result, err := c.GetKubernetesClusterRole(context.Background(), tt.identity)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestDeleteKubernetesClusterRole(t *testing.T) {
	tests := []struct {
		name         string
		identity     string
		serverStatus int
		expectError  bool
	}{
		{
			name:         "successful deletion",
			identity:     "role-1",
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "not found",
			identity:     "non-existent",
			serverStatus: http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, KubernetesClusterRoleEndpoint+"/"):
					w.WriteHeader(tt.serverStatus)
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			err = c.DeleteClusterRole(context.Background(), tt.identity)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddRoleRule(t *testing.T) {
	tests := []struct {
		name           string
		identity       string
		rule           AddKubernetesClusterRolePermissionRule
		serverResponse *KubernetesClusterRolePermissionRule
		serverStatus   int
		expectError    bool
	}{
		{
			name:     "successful rule addition",
			identity: "role-1",
			rule: AddKubernetesClusterRolePermissionRule{
				Resources: []string{"pods"},
				Verbs:     []KubernetesClusterRolePermissionVerb{KubernetesClusterRolePermissionVerbGet},
				ApiGroups: []string{""},
			},
			serverResponse: &KubernetesClusterRolePermissionRule{
				Identity:  "rule-1",
				Resources: []string{"pods"},
				Verbs:     []KubernetesClusterRolePermissionVerb{KubernetesClusterRolePermissionVerbGet},
				ApiGroups: []string{""},
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:     "server error",
			identity: "role-1",
			rule: AddKubernetesClusterRolePermissionRule{
				Resources: []string{"pods"},
				Verbs:     []KubernetesClusterRolePermissionVerb{KubernetesClusterRolePermissionVerbGet},
			},
			serverResponse: nil,
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" && strings.Contains(r.URL.Path, "/rules") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					if tt.serverStatus == http.StatusOK && tt.serverResponse != nil {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			result, err := c.AddClusterRoleRule(context.Background(), tt.identity, tt.rule)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestDeleteRuleFromRole(t *testing.T) {
	tests := []struct {
		name         string
		identity     string
		ruleIdentity string
		serverStatus int
		expectError  bool
	}{
		{
			name:         "successful rule deletion",
			identity:     "role-1",
			ruleIdentity: "rule-1",
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:         "rule not found",
			identity:     "role-1",
			ruleIdentity: "non-existent",
			serverStatus: http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "DELETE" && strings.Contains(r.URL.Path, "/rules/") {
					w.WriteHeader(tt.serverStatus)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			err = c.DeleteClusterRoleRule(context.Background(), tt.identity, tt.ruleIdentity)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListRoleBindings(t *testing.T) {
	tests := []struct {
		name           string
		identity       string
		serverResponse []KubernetesClusterRoleBinding
		serverStatus   int
		expectError    bool
		expectedCount  int
	}{
		{
			name:     "successful list",
			identity: "role-1",
			serverResponse: []KubernetesClusterRoleBinding{
				{
					Identity: "binding-1",
					Name:     "test-binding",
					Slug:     "test-binding",
				},
			},
			serverStatus:  http.StatusOK,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:           "empty list",
			identity:       "role-1",
			serverResponse: []KubernetesClusterRoleBinding{},
			serverStatus:   http.StatusOK,
			expectError:    false,
			expectedCount:  0,
		},
		{
			name:           "server error",
			identity:       "role-1",
			serverResponse: []KubernetesClusterRoleBinding{},
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/bindings") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					if tt.serverStatus == http.StatusOK {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			result, err := c.ListClusterRoleBindings(context.Background(), tt.identity)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}
}

func TestCreateRoleBinding(t *testing.T) {
	tests := []struct {
		name           string
		identity       string
		createRequest  CreateKubernetesClusterRoleBinding
		serverResponse *KubernetesClusterRoleBinding
		serverStatus   int
		expectError    bool
	}{
		{
			name:     "successful binding creation with user",
			identity: "role-1",
			createRequest: CreateKubernetesClusterRoleBinding{
				Name:         "test-binding",
				Description:  "Test binding",
				UserIdentity: stringPtr("user-1"),
			},
			serverResponse: &KubernetesClusterRoleBinding{
				Identity: "binding-1",
				Name:     "test-binding",
				Slug:     "test-binding",
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:     "successful binding creation with team",
			identity: "role-1",
			createRequest: CreateKubernetesClusterRoleBinding{
				Name:         "test-binding",
				Description:  "Test binding",
				TeamIdentity: stringPtr("team-1"),
			},
			serverResponse: &KubernetesClusterRoleBinding{
				Identity: "binding-1",
				Name:     "test-binding",
				Slug:     "test-binding",
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name:     "server error",
			identity: "role-1",
			createRequest: CreateKubernetesClusterRoleBinding{
				Name:         "test-binding",
				Description:  "Test binding",
				UserIdentity: stringPtr("user-1"),
			},
			serverResponse: nil,
			serverStatus:   http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/bindings") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.serverStatus)
					if tt.serverStatus == http.StatusOK && tt.serverResponse != nil {
						json.NewEncoder(w).Encode(tt.serverResponse)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			result, err := c.CreateClusterRoleBinding(context.Background(), tt.identity, tt.createRequest)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestDeleteRoleBinding(t *testing.T) {
	tests := []struct {
		name                string
		identity            string
		roleBindingIdentity string
		serverStatus        int
		expectError         bool
	}{
		{
			name:                "successful binding deletion",
			identity:            "role-1",
			roleBindingIdentity: "binding-1",
			serverStatus:        http.StatusOK,
			expectError:         false,
		},
		{
			name:                "binding not found",
			identity:            "role-1",
			roleBindingIdentity: "non-existent",
			serverStatus:        http.StatusNotFound,
			expectError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "DELETE" && strings.Contains(r.URL.Path, "/bindings/") {
					w.WriteHeader(tt.serverStatus)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			baseClient, err := client.NewClient(client.WithBaseURL(server.URL))
			require.NoError(t, err)
			c, err := New(baseClient)
			require.NoError(t, err)

			err = c.DeleteClusterRoleBinding(context.Background(), tt.identity, tt.roleBindingIdentity)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
