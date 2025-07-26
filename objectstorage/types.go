package objectstorage

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type ObjectStorageBucket struct {
	Identity     string             `json:"identity"`
	Organisation *base.Organisation `json:"organisation,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`

	// Name is the name of the bucket
	Name string `json:"name"`
	// Policy is the policy of the bucket
	Policy PolicyDocument `json:"policy"`
	// Public is a flag that indicates if the bucket is public
	Public bool `json:"public"`
	// Status is the status of the bucket
	Status string `json:"status"`
	// Usage is the usage of the bucket
	Usage Usage `json:"usage"`

	// Endpoint for the bucket, is collected from the CR
	Endpoint string `json:"endpoint"`

	// Region is the region of the bucket
	Region *iaas.Region `json:"cloudRegion,omitempty"`
}

type Usage struct {
	TotalSizeGB  int `json:"total_size_gb"`
	TotalObjects int `json:"total_objects"`
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
}

type UpdateBucketRequest struct {
	// Public is a flag that indicates if the bucket can be accessed by the public.
	Public bool `json:"public"`

	// PolicyDocument is the policy document for the bucket.
	PolicyDocument *PolicyDocument `json:"policy,omitempty"`
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
	AWS interface{} `json:"AWS"` // can be string or []string
}
