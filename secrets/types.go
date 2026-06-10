package secrets

import (
	"time"

	"github.com/thalassa-cloud/client-go/kms"
)

type SecretVersion struct {
	Version     int        `json:"version"`
	Status      string     `json:"status,omitempty"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	DestroyedAt *time.Time `json:"destroyedAt,omitempty"`
}

type SecretPolicyStatement struct {
	Effect     string   `json:"effect"`
	Actions    []string `json:"actions,omitempty"`
	Principals []string `json:"principals,omitempty"`
}

type SecretPolicy struct {
	Statements []SecretPolicyStatement `json:"statements,omitempty"`
}

type GenerateSecret struct {
	Length       int    `json:"length,omitempty"`
	CharacterSet string `json:"characterSet,omitempty"`
}

type Secret struct {
	Path           string            `json:"path"`
	Description    string            `json:"description,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	KmsKey         *kms.KmsKey       `json:"kmsKey,omitempty"`
	KmsKeyIdentity string            `json:"kmsKeyIdentity,omitempty"`
	CurrentVersion int               `json:"currentVersion"`
	LastAccessedAt *time.Time        `json:"lastAccessedAt,omitempty"`
	AccessPolicy   *SecretPolicy     `json:"accessPolicy,omitempty"`
	Versions       []SecretVersion   `json:"versions,omitempty"`
	CreatedAt      time.Time         `json:"createdAt,omitempty"`
	UpdatedAt      time.Time         `json:"updatedAt,omitempty"`
}

type BrowseSecretsResponse struct {
	Path     string   `json:"path"`
	Prefixes []string `json:"prefixes,omitempty"`
	Secrets  []Secret `json:"secrets,omitempty"`
}

type CreateSecretRequest struct {
	Path            string            `json:"path"`
	Description     string            `json:"description,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	KmsKeyIdentity  string            `json:"kmsKeyIdentity"`
	SecretString    string            `json:"secretString,omitempty"`
	SecretKeyValues map[string]string `json:"secretKeyValues,omitempty"`
	GenerateSecret  *GenerateSecret   `json:"generateSecret,omitempty"`
	AccessPolicy    *SecretPolicy     `json:"accessPolicy,omitempty"`
}

type PutSecretValueRequest struct {
	SecretString    string            `json:"secretString,omitempty"`
	SecretKeyValues map[string]string `json:"secretKeyValues,omitempty"`
	GenerateSecret  *GenerateSecret   `json:"generateSecret,omitempty"`
}

type PutSecretValueResponse struct {
	Path           string `json:"path"`
	Version        int    `json:"version"`
	KmsKeyIdentity string `json:"kmsKeyIdentity,omitempty"`
	KmsKeyVersion  string `json:"kmsKeyVersion,omitempty"`
}

type GetSecretValueRequest struct {
	Version *int `json:"version,omitempty"`
}

type GetSecretValueResponse struct {
	Path            string            `json:"path"`
	Version         int               `json:"version"`
	SecretString    string            `json:"secretString,omitempty"`
	SecretKeyValues map[string]string `json:"secretKeyValues,omitempty"`
	KmsKeyIdentity  string            `json:"kmsKeyIdentity"`
	KmsKeyVersion   string            `json:"kmsKeyVersion"`
}

type UpdateAccessPolicyRequest struct {
	AccessPolicy SecretPolicy `json:"accessPolicy"`
}
