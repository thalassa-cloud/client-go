package iam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func newTestIAMClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	c, err := client.NewClient(
		client.WithBaseURL(serverURL),
		client.WithAuthCustom(),
		client.WithOrganisation("acme"),
		client.WithProject("prj-test"),
	)
	require.NoError(t, err)
	iamClient, err := New(c)
	require.NoError(t, err)
	return iamClient
}

func assertIAMScopedHeaders(t *testing.T, r *http.Request) {
	t.Helper()
	assert.Equal(t, "acme", r.Header.Get("X-Organisation-Identity"))
	assert.Equal(t, "prj-test", r.Header.Get("X-Project-Identity"))
}

func TestListIamPolicies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies", r.URL.Path)
		assertIAMScopedHeaders(t, r)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]IamPolicy{
			{Identity: "pol-1", Name: "admins", Slug: "admins"},
		}))
	}))
	defer server.Close()

	policies, err := newTestIAMClient(t, server.URL).ListIamPolicies(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, policies, 1)
	assert.Equal(t, "pol-1", policies[0].Identity)
}

func TestGetIamPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1", r.URL.Path)
		assertIAMScopedHeaders(t, r)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamPolicy{
			Identity: "pol-1",
			Name:     "admins",
			Rules: []IamPolicyPermissionRule{
				{Identity: "rule-1", Resources: []string{"vpc"}, Permissions: []PermissionType{PermissionTypeRead}},
			},
		}))
	}))
	defer server.Close()

	policy, err := newTestIAMClient(t, server.URL).GetIamPolicy(context.Background(), "pol-1")
	require.NoError(t, err)
	assert.Equal(t, "pol-1", policy.Identity)
	require.Len(t, policy.Rules, 1)
}

func TestGetIamPolicyValidation(t *testing.T) {
	_, err := newTestIAMClient(t, "http://example.com").GetIamPolicy(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "identity is required")
}

func TestGetIamPolicyNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"message": "policy not found"}))
	}))
	defer server.Close()

	_, err := newTestIAMClient(t, server.URL).GetIamPolicy(context.Background(), "missing")
	require.Error(t, err)
	assert.True(t, client.IsNotFound(err))
}

func TestCreateIamPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies", r.URL.Path)
		assertIAMScopedHeaders(t, r)

		var body CreateIamPolicyRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "admins", body.Name)
		assert.True(t, body.ReplicateToChildren)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(IamPolicy{
			Identity:            "pol-1",
			Name:                body.Name,
			Slug:                "admins",
			ReplicateToChildren: body.ReplicateToChildren,
			CreatedAt:           time.Now().UTC(),
		}))
	}))
	defer server.Close()

	policy, err := newTestIAMClient(t, server.URL).CreateIamPolicy(context.Background(), CreateIamPolicyRequest{
		Name:                "admins",
		ReplicateToChildren: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "pol-1", policy.Identity)
}

func TestCreateIamPolicyValidation(t *testing.T) {
	_, err := newTestIAMClient(t, "http://example.com").CreateIamPolicy(context.Background(), CreateIamPolicyRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestUpdateIamPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1", r.URL.Path)

		var body UpdateIamPolicyRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "updated", body.Description)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamPolicy{Identity: "pol-1", Description: body.Description}))
	}))
	defer server.Close()

	policy, err := newTestIAMClient(t, server.URL).UpdateIamPolicy(context.Background(), "pol-1", UpdateIamPolicyRequest{
		Description: "updated",
	})
	require.NoError(t, err)
	assert.Equal(t, "updated", policy.Description)
}

func TestDeleteIamPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := newTestIAMClient(t, server.URL).DeleteIamPolicy(context.Background(), "pol-1")
	require.NoError(t, err)
}

func TestListResourceTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/policies/resources", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]string{"vpc", "kubernetes_cluster"}))
	}))
	defer server.Close()

	types, err := newTestIAMClient(t, server.URL).ListResourceTypes(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"vpc", "kubernetes_cluster"}, types)
}

func TestListProjectMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/policies/members", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]ProjectMember{
			{Bindings: []ProjectMemberPolicyAccess{{Source: ProjectMemberSourceIamPolicy, PolicyIdentity: "pol-1"}}},
		}))
	}))
	defer server.Close()

	members, err := newTestIAMClient(t, server.URL).ListProjectMembers(context.Background())
	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, "pol-1", members[0].Bindings[0].PolicyIdentity)
}

func TestAddIamPolicyRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1/rules", r.URL.Path)

		var body AddIamPolicyRuleRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []string{"vpc"}, body.Resources)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(IamPolicyPermissionRule{
			Identity:    "rule-1",
			Resources:   body.Resources,
			Permissions: body.Permissions,
		}))
	}))
	defer server.Close()

	rule, err := newTestIAMClient(t, server.URL).AddIamPolicyRule(context.Background(), "pol-1", AddIamPolicyRuleRequest{
		Resources:   []string{"vpc"},
		Permissions: []PermissionType{PermissionTypeRead},
	})
	require.NoError(t, err)
	assert.Equal(t, "rule-1", rule.Identity)
}

func TestAddIamPolicyRuleValidation(t *testing.T) {
	c := newTestIAMClient(t, "http://example.com")
	_, err := c.AddIamPolicyRule(context.Background(), "pol-1", AddIamPolicyRuleRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resources is required")
}

func TestDeleteIamPolicyRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1/rules/rule-1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := newTestIAMClient(t, server.URL).DeleteIamPolicyRule(context.Background(), "pol-1", "rule-1")
	require.NoError(t, err)
}

func TestCreateIamPolicyBinding(t *testing.T) {
	userIdentity := "usr-1"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1/bindings", r.URL.Path)

		var body CreateIamPolicyBindingRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "binding", body.Name)
		require.NotNil(t, body.UserIdentity)
		assert.Equal(t, userIdentity, *body.UserIdentity)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(IamPolicyBinding{Identity: "bind-1", Name: body.Name}))
	}))
	defer server.Close()

	binding, err := newTestIAMClient(t, server.URL).CreateIamPolicyBinding(context.Background(), "pol-1", CreateIamPolicyBindingRequest{
		Name:         "binding",
		UserIdentity: &userIdentity,
	})
	require.NoError(t, err)
	assert.Equal(t, "bind-1", binding.Identity)
}

func TestCreateIamPolicyBindingValidation(t *testing.T) {
	c := newTestIAMClient(t, "http://example.com")
	_, err := c.CreateIamPolicyBinding(context.Background(), "pol-1", CreateIamPolicyBindingRequest{Name: "x"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one of userIdentity or serviceAccountIdentity")
}

func TestListIamPolicyBindings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/projects/iam/policies/pol-1/bindings", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]IamPolicyBinding{{Identity: "bind-1"}}))
	}))
	defer server.Close()

	bindings, err := newTestIAMClient(t, server.URL).ListIamPolicyBindings(context.Background(), "pol-1", nil)
	require.NoError(t, err)
	require.Len(t, bindings, 1)
}

func TestUpdateIamPolicyBinding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1/bindings/bind-1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(IamPolicyBinding{Identity: "bind-1", Description: "d"}))
	}))
	defer server.Close()

	binding, err := newTestIAMClient(t, server.URL).UpdateIamPolicyBinding(context.Background(), "pol-1", "bind-1", UpdateIamPolicyBindingRequest{
		Description: "d",
	})
	require.NoError(t, err)
	assert.Equal(t, "d", binding.Description)
}

func TestDeleteIamPolicyBinding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/projects/iam/policies/pol-1/bindings/bind-1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := newTestIAMClient(t, server.URL).DeleteIamPolicyBinding(context.Background(), "pol-1", "bind-1")
	require.NoError(t, err)
}
