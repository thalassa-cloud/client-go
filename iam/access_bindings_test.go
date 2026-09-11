package iam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListResourceAccessBindings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects/iam/access-bindings", r.URL.Path)
		assertIAMScopedHeaders(t, r)

		assert.Equal(t, []string{"vpc", "subnet"}, r.URL.Query()["resourceType"])
		assert.Equal(t, "vpc-1", r.URL.Query().Get("resourceIdentity"))
		assert.Equal(t, []string{"read", "list"}, r.URL.Query()["permissions"])
		assert.Equal(t, []string{"iam_policy"}, r.URL.Query()["bindingTypes"])

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(ResourceAccessBindingsResponse{
			ResourceIdentity: "vpc-1",
			Bindings: []ResourceAccessBindingItem{
				{
					ResourceType:      "vpc",
					BindingType:       "iam_policy",
					BindingIdentity:   "bind-1",
					PolicyIdentity:    "pol-1",
					PrincipalType:     ResourceAccessBindingPrincipalUser,
					PrincipalIdentity: "usr-1",
					Permissions:       []string{"read"},
				},
			},
		}))
	}))
	defer server.Close()

	result, err := newTestIAMClient(t, server.URL).ListResourceAccessBindings(context.Background(), ListResourceAccessBindingsRequest{
		ResourceTypes:    []string{"vpc", "subnet"},
		ResourceIdentity: "vpc-1",
		Permissions:      []PermissionType{PermissionTypeRead, PermissionTypeList},
		BindingTypes:     []string{"iam_policy"},
	})
	require.NoError(t, err)
	assert.Equal(t, "vpc-1", result.ResourceIdentity)
	require.Len(t, result.Bindings, 1)
	assert.Equal(t, "bind-1", result.Bindings[0].BindingIdentity)
}

func TestListResourceAccessBindingsValidation(t *testing.T) {
	c := newTestIAMClient(t, "http://example.com")

	_, err := c.ListResourceAccessBindings(context.Background(), ListResourceAccessBindingsRequest{
		ResourceTypes: []string{"vpc"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resourceIdentity is required")

	_, err = c.ListResourceAccessBindings(context.Background(), ListResourceAccessBindingsRequest{
		ResourceIdentity: "vpc-1",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resourceType is required")
}
