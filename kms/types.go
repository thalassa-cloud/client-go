package kms

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

type KmsKeyStatus string

const (
	KmsKeyStatusActive          KmsKeyStatus = "active"
	KmsKeyStatusDisabled        KmsKeyStatus = "disabled"
	KmsKeyStatusPendingDeletion KmsKeyStatus = "pending_deletion"
)

type KmsKeyType string

const (
	KmsKeyTypeAES256GCM96 KmsKeyType = "aes256-gcm96"
	KmsKeyTypeRSA2048     KmsKeyType = "rsa-2048"
	KmsKeyTypeRSA3072     KmsKeyType = "rsa-3072"
	KmsKeyTypeRSA4096     KmsKeyType = "rsa-4096"
	KmsKeyTypeECP256      KmsKeyType = "ec-p256"
	KmsKeyTypeECP384      KmsKeyType = "ec-p384"
	KmsKeyTypeECP521      KmsKeyType = "ec-p521"
	KmsKeyTypeEd25519     KmsKeyType = "ed25519"
	KmsKeyTypeHMACSHA256  KmsKeyType = "hmac-sha256"
	KmsKeyTypeHMACSHA384  KmsKeyType = "hmac-sha384"
	KmsKeyTypeHMACSHA512  KmsKeyType = "hmac-sha512"
)

type KmsKeyVersion struct {
	Version   int       `json:"version"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type KmsKey struct {
	Identity             string             `json:"identity"`
	Name                 string             `json:"name"`
	Slug                 string             `json:"slug"`
	Description          string             `json:"description,omitempty"`
	Labels               map[string]string  `json:"labels,omitempty"`
	Annotations          map[string]string  `json:"annotations,omitempty"`
	ProjectID            string             `json:"projectId,omitempty"`
	KeyType              KmsKeyType         `json:"keyType"`
	Status               KmsKeyStatus       `json:"status"`
	ExportAllowed        bool               `json:"exportAllowed"`
	Imported             bool               `json:"imported,omitempty"`
	KeyRotationEnabled   bool               `json:"keyRotationEnabled"`
	RotationPeriodInDays *int               `json:"rotationPeriodInDays,omitempty"`
	LatestVersion        int                `json:"latestVersion,omitempty"`
	Versions             []KmsKeyVersion    `json:"versions,omitempty"`
	CreatedAt            time.Time          `json:"createdAt,omitempty"`
	UpdatedAt            time.Time          `json:"updatedAt,omitempty"`
	ObjectVersion        int                `json:"objectVersion,omitempty"`
	Organisation         *base.Organisation `json:"organisation,omitempty"`
}

type KmsSummaryRegion struct {
	Identity         string `json:"identity"`
	Name             string `json:"name,omitempty"`
	Slug             string `json:"slug"`
	KmsAvailable     bool   `json:"kmsAvailable"`
	SecretsAvailable bool   `json:"secretsAvailable,omitempty"`
}

type KmsSummary struct {
	FeatureEnabled bool               `json:"featureEnabled"`
	Regions        []KmsSummaryRegion `json:"regions,omitempty"`
}

type ListKeysRequest struct {
	Filters []ListKeysFilter
}

type ListKeysFilter interface {
	ToParams() map[string]string
}

type CreateKmsKeyRequest struct {
	Name                 string            `json:"name"`
	Description          string            `json:"description,omitempty"`
	Labels               map[string]string `json:"labels,omitempty"`
	Annotations          map[string]string `json:"annotations,omitempty"`
	KeyType              KmsKeyType        `json:"keyType,omitempty"`
	ExportAllowed        bool              `json:"exportAllowed,omitempty"`
	KeyRotationEnabled   bool              `json:"keyRotationEnabled,omitempty"`
	RotationPeriodInDays *int              `json:"rotationPeriodInDays,omitempty"`
	ImportKeyMaterial    string            `json:"importKeyMaterial,omitempty"`
	HashFunction         string            `json:"hashFunction,omitempty"`
	AllowRotation        bool              `json:"allowRotation,omitempty"`
}

type UpdateRotationRequest struct {
	KeyRotationEnabled   *bool `json:"keyRotationEnabled,omitempty"`
	RotationPeriodInDays *int  `json:"rotationPeriodInDays,omitempty"`
}

type EncryptRequest struct {
	Plaintext  string `json:"plaintext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type EncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type DecryptRequest struct {
	Ciphertext string `json:"ciphertext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type DecryptResponse struct {
	Plaintext  string `json:"plaintext"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type SignRequest struct {
	Message    string `json:"message"`
	KeyVersion string `json:"keyVersion,omitempty"`
	Hash       string `json:"hash,omitempty"`
}

type SignResponse struct {
	Signature  string `json:"signature"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type VerifySignatureRequest struct {
	Message    string `json:"message"`
	Signature  string `json:"signature"`
	Hash       string `json:"hash,omitempty"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type VerifySignatureResponse struct {
	Valid      bool   `json:"valid"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type HMACRequest struct {
	Message    string `json:"message"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type HMACResponse struct {
	MAC        string `json:"mac"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type VerifyHMACRequest struct {
	Message    string `json:"message"`
	MAC        string `json:"mac"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type VerifyHMACResponse struct {
	Valid      bool   `json:"valid"`
	KeyVersion string `json:"keyVersion,omitempty"`
}

type GetPublicKeyResponse struct {
	PublicKey  string `json:"publicKey"`
	KeyVersion string `json:"keyVersion,omitempty"`
	KeyType    string `json:"keyType,omitempty"`
}

type ExportKeyRequest struct {
	KeyVersion string `json:"keyVersion,omitempty"`
}

type ExportKeyResponse struct {
	KeyMaterial string `json:"keyMaterial"`
	KeyVersion  string `json:"keyVersion,omitempty"`
}

type WrappingKeyResponse struct {
	PublicKey string `json:"publicKey"`
	Algorithm string `json:"algorithm,omitempty"`
}
