package objectstorage

import (
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type ObjectStorageBucketVersioning string

const (
	ObjectStorageBucketVersioningDisabled  ObjectStorageBucketVersioning = "Disabled"
	ObjectStorageBucketVersioningEnabled   ObjectStorageBucketVersioning = "Enabled"
	ObjectStorageBucketVersioningSuspended ObjectStorageBucketVersioning = "Suspended"
)

type ObjectStorageBucket struct {
	Identity     string             `json:"identity"`
	Organisation *base.Organisation `json:"organisation,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`

	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Name is the name of the bucket
	Name string `json:"name"`
	// Policy is the policy of the bucket
	Policy PolicyDocument `json:"policy"`
	// Public is a flag that indicates if the bucket is public
	Public bool `json:"public"`
	// Status is the status of the bucket
	Status string `json:"status"`
	// Endpoint for the bucket
	Endpoint string `json:"endpoint"`
	// Usage is the usage of the bucket
	Usage ObjectStorageBucketUsage `json:"usage"`
	// Versioning is the versioning of the bucket
	Versioning ObjectStorageBucketVersioning `json:"versioning"`
	// ObjectLockEnabled is the object lock of the bucket
	ObjectLockEnabled bool `json:"objectLockEnabled"`
	// Region is the region of the bucket
	Region *iaas.Region `json:"cloudRegion,omitempty"`
}

type ObjectStorageBucketUsage struct {
	TotalSizeGB  float64 `json:"total_size_gb"`
	TotalObjects int64   `json:"total_objects"`
}

type CreateBucketRequest struct {
	// BucketName is the name of the bucket.
	BucketName string `json:"bucketName"`
	// Public is a flag that indicates if the bucket can be accessed by the public.
	// When set to false, it blocks all public access to the bucket.
	Public bool `json:"public"`
	// Region is the region of the bucket.
	Region string `json:"region"`
	// PolicyDocument is the policy document for the bucket.
	PolicyDocument *PolicyDocument `json:"policy,omitempty"`
	// Labels is the labels of the bucket.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is the annotations of the bucket.
	Annotations map[string]string `json:"annotations,omitempty"`
	// Versioning is the versioning of the bucket.
	Versioning ObjectStorageBucketVersioning `json:"versioning"`
	// ObjectLockEnabled is the object lock enabled of the bucket.
	ObjectLockEnabled bool `json:"objectLockEnabled"`
}

type UpdateBucketRequest struct {
	// Public is a flag that indicates if the bucket can be accessed by the public.
	Public bool `json:"public"`
	// PolicyDocument is the policy document for the bucket.
	PolicyDocument *PolicyDocument `json:"policy,omitempty"`
	// Versioning is the versioning of the bucket.
	Versioning ObjectStorageBucketVersioning `json:"versioning"`
	// ObjectLockEnabled is the object lock enabled of the bucket.
	ObjectLockEnabled *bool `json:"objectLockEnabled"`

	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// PolicyDocument represents a full S3 bucket policy.
type PolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// Statement defines an individual rule in the policy.
type Statement struct {
	Sid       string      `json:"Sid,omitempty"`
	Effect    string      `json:"Effect"`
	Principal Principal   `json:"Principal"`
	Action    interface{} `json:"Action"` // can be string or []string
	Resource  []string    `json:"Resource"`
	Condition interface{} `json:"Condition,omitempty"`
}

// Principal defines which user(s) the statement applies to.
type Principal struct {
	AWS      interface{} `json:"AWS,omitempty"`      // can be string or []string
	Thalassa interface{} `json:"Thalassa,omitempty"` // can be string or []string
}

type PrincipalARN string

func (p PrincipalARN) Validate() error {
	if p == "" {
		return fmt.Errorf("principal ARN is required")
	}
	return nil
}
