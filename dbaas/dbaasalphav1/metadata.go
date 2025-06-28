package dbaasalphav1

// Labels are key-value pairs that can be used to label resources.
// Keys and values must be RFC 1123 compliant.
// Label keys and values must be 1-63 characters long, and must conform to the following
// - contain at most 63 characters
// - contain only lowercase alphanumeric characters or '-'
// - start with an alphanumeric character
// - end with an alphanumeric character
type Labels map[string]string

// Annotations are key-value pairs that can be used to annotate resources.
// Keys must be RFC 1123 compliant, but the values may contain any ascii characters.
type Annotations map[string]string

type ObjectStatus string

const (
	ObjectStatusCreating ObjectStatus = "creating"
	ObjectStatusReady    ObjectStatus = "ready"
	ObjectStatusDeleting ObjectStatus = "deleting"
	ObjectStatusUpdating ObjectStatus = "updating"
	ObjectStatusDeleted  ObjectStatus = "deleted"
	ObjectStatusFailed   ObjectStatus = "failed"
)
