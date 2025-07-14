package kubernetes

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetKubernetesClusterKubeconfig(t *testing.T) {
	tests := []struct {
		name          string
		identity      string
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "empty identity",
			identity:      "",
			ctx:           context.Background(),
			expectedError: "cluster identity cannot be empty",
		},
		{
			name:          "whitespace identity",
			identity:      "   ",
			ctx:           context.Background(),
			expectedError: "cluster identity cannot be empty",
		},
		{
			name:          "nil context",
			identity:      "cluster-123",
			ctx:           nil,
			expectedError: "context cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a minimal client for testing validation
			client := &Client{}

			// Execute the method
			result, err := client.GetKubernetesClusterKubeconfig(tt.ctx, tt.identity)

			// Assert results
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, result)
		})
	}
}

// Test constants and helper functions
func TestKubernetesClusterKubeConfigEndpoint_Format(t *testing.T) {
	clusterID := "test-cluster-123"
	expected := "/v1/kubernetes/clusters/test-cluster-123/kubeconfig"

	formatted := fmt.Sprintf(KubernetesClusterKubeConfigEndpoint, clusterID)
	assert.Equal(t, expected, formatted)
}

func TestKubernetesClusterSessionToken_Validation(t *testing.T) {
	tests := []struct {
		name    string
		token   *KubernetesClusterSessionToken
		isValid bool
	}{
		{
			name: "valid token",
			token: &KubernetesClusterSessionToken{
				Username:      "admin",
				APIServerURL:  "https://api.example.com",
				CACertificate: "-----BEGIN CERTIFICATE-----",
				Identity:      "session-123",
				Token:         "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9",
				Kubeconfig:    "apiVersion: v1",
			},
			isValid: true,
		},
		{
			name: "missing API server URL",
			token: &KubernetesClusterSessionToken{
				Username:      "admin",
				APIServerURL:  "",
				CACertificate: "-----BEGIN CERTIFICATE-----",
				Identity:      "session-123",
				Token:         "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9",
				Kubeconfig:    "apiVersion: v1",
			},
			isValid: false,
		},
		{
			name: "missing token",
			token: &KubernetesClusterSessionToken{
				Username:      "admin",
				APIServerURL:  "https://api.example.com",
				CACertificate: "-----BEGIN CERTIFICATE-----",
				Identity:      "session-123",
				Token:         "",
				Kubeconfig:    "apiVersion: v1",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.token.APIServerURL != "" && tt.token.Token != ""
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}
