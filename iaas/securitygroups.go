package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	SecurityGroupEndpoint = "/v1/security-groups"
)

type CreateSecurityGroupRequest struct {
	// Name is the name of the security group
	Name string `json:"name" validate:"required,min=1,max=16,notblank,ascii,is-trimmed"`
	// Description is the description of the security group
	Description string `json:"description" validate:"omitempty,max=255"`
	// Labels are the labels of the security group
	Labels Labels `json:"labels"`
	// Annotations are the annotations of the security group
	Annotations Annotations `json:"annotations"`
	// VpcIdentity is the identity of the VPC that the security group belongs to
	VpcIdentity string `json:"vpcIdentity" validate:"required"`
	// AllowSameGroupTraffic is a flag that indicates if the security group allows traffic between instances in the same security group
	AllowSameGroupTraffic bool `json:"allowSameGroupTraffic"`
	// IngressRules are the ingress rules of the security group
	IngressRules []SecurityGroupRule `json:"ingressRules" validate:"omitempty,dive"`
	// EgressRules are the egress rules of the security group
	EgressRules []SecurityGroupRule `json:"egressRules" validate:"omitempty,dive"`
}

type UpdateSecurityGroupRequest struct {
	// Name is the name of the security group
	Name string `json:"name" validate:"min=1,max=16,notblank,ascii,is-trimmed"`
	// Description is the description of the security group
	Description string `json:"description" validate:"omitempty,max=255"`
	// Labels are the labels of the security group
	Labels Labels `json:"labels"`
	// Annotations are the annotations of the security group
	Annotations Annotations `json:"annotations"`
	// ObjectVersion is the version of the security group
	ObjectVersion int `json:"objectVersion"`
	// AllowSameGroupTraffic is a flag that indicates if the security group allows traffic between instances in the same security group
	AllowSameGroupTraffic bool `json:"allowSameGroupTraffic"`
	// IngressRules are the ingress rules of the security group
	IngressRules []SecurityGroupRule `json:"ingressRules" validate:"omitempty,dive"`
	// EgressRules are the egress rules of the security group
	EgressRules []SecurityGroupRule `json:"egressRules" validate:"omitempty,dive"`
}

type SecurityGroupStatus string

const (
	SecurityGroupStatusProvisioning SecurityGroupStatus = "provisioning"
	SecurityGroupStatusActive       SecurityGroupStatus = "active"
	SecurityGroupStatusReady        SecurityGroupStatus = "ready"
	SecurityGroupStatusDeleting     SecurityGroupStatus = "deleting"
	SecurityGroupStatusError        SecurityGroupStatus = "error"
)

type SecurityGroup struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`

	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	// Organisation is the organisation that the security group belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Vpc is the VPC that the security group belongs to
	Vpc *Vpc `json:"vpc,omitempty"`

	// Status is the status of the security group
	Status SecurityGroupStatus `json:"status,omitempty"`
	// AllowSameGroupTraffic is a flag that indicates if the security group allows traffic between instances in the same security group
	AllowSameGroupTraffic bool `json:"allowSameGroupTraffic" validate:"required"`
	// IngressRules are the ingress rules of the security group
	IngressRules []SecurityGroupRule `json:"ingressRules" validate:"required,dive,max=60,min=0"`
	// EgressRules are the egress rules of the security group
	EgressRules []SecurityGroupRule `json:"egressRules" validate:"required,dive,max=60,min=0"`
}

type SecurityGroupRule struct {
	// Name is the name of the rule
	Name string `json:"name" validate:"omitempty,max=255"`
	// IPVersion is the IP version of the rule
	IPVersion SecurityGroupIPVersion `json:"ipVersion" validate:"required"`
	// Protocol is the protocol of the rule
	Protocol SecurityGroupRuleProtocol `json:"protocol" validate:"required"`
	// Priority is the priority of the rule
	Priority int32 `json:"priority" validate:"required,gte=1,lte=200"`
	// RemoteType is the type of the remote address
	RemoteType SecurityGroupRuleRemoteType `json:"remoteType" validate:"required"`
	// RemoteAddress is the IP address or CIDR block that the rule applies to
	RemoteAddress *string `json:"remoteAddress" validate:"omitempty"`
	// RemoteSecurityGroupIdentity is the identity of the security group that the rule applies to
	RemoteSecurityGroupIdentity *string `json:"remoteSecurityGroupIdentity" validate:"omitempty"`
	// PortRangeMin is the minimum port of the rule
	PortRangeMin int32 `json:"portRangeMin" validate:"omitempty,gte=1,lte=65535"`
	// PortRangeMax is the maximum port of the rule
	PortRangeMax int32 `json:"portRangeMax" validate:"omitempty,gte=0,lte=65535"`
	// Policy is the policy of the rule
	Policy SecurityGroupRulePolicy `json:"policy" validate:"required"`
}

type SecurityGroupRuleProtocol string

const (
	SecurityGroupRuleProtocolAll  SecurityGroupRuleProtocol = "all"
	SecurityGroupRuleProtocolTCP  SecurityGroupRuleProtocol = "tcp"
	SecurityGroupRuleProtocolUDP  SecurityGroupRuleProtocol = "udp"
	SecurityGroupRuleProtocolICMP SecurityGroupRuleProtocol = "icmp"
)

type SecurityGroupRulePolicy string

const (
	SecurityGroupRulePolicyAllow SecurityGroupRulePolicy = "allow"
	SecurityGroupRulePolicyDrop  SecurityGroupRulePolicy = "drop"
)

type SecurityGroupRuleRemoteType string

const (
	SecurityGroupRuleRemoteTypeAddress       SecurityGroupRuleRemoteType = "address"
	SecurityGroupRuleRemoteTypeSecurityGroup SecurityGroupRuleRemoteType = "securityGroup"
)

type SecurityGroupIPVersion string

const (
	SecurityGroupIPVersionIPv4 SecurityGroupIPVersion = "ipv4"
	SecurityGroupIPVersionIPv6 SecurityGroupIPVersion = "ipv6"
)

type ListSecurityGroupsRequest struct {
	Filters []filters.Filter
}

// ListSecurityGroups lists all security groups for a given organisation.
func (c *Client) ListSecurityGroups(ctx context.Context, listRequest *ListSecurityGroupsRequest) ([]SecurityGroup, error) {
	securityGroups := []SecurityGroup{}
	req := c.R().SetResult(&securityGroups)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, SecurityGroupEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return securityGroups, err
	}
	return securityGroups, nil
}

// GetSecurityGroup retrieves a specific security group by its identity.
func (c *Client) GetSecurityGroup(ctx context.Context, identity string) (*SecurityGroup, error) {
	var securityGroup *SecurityGroup
	req := c.R().SetResult(&securityGroup)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", SecurityGroupEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return securityGroup, err
	}
	return securityGroup, nil
}

// CreateSecurityGroup creates a new security group.
func (c *Client) CreateSecurityGroup(ctx context.Context, create CreateSecurityGroupRequest) (*SecurityGroup, error) {
	var securityGroup *SecurityGroup
	req := c.R().
		SetBody(create).SetResult(&securityGroup)

	resp, err := c.Do(ctx, req, client.POST, SecurityGroupEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return securityGroup, err
	}
	return securityGroup, nil
}

// UpdateSecurityGroup updates an existing security group.
func (c *Client) UpdateSecurityGroup(ctx context.Context, identity string, update UpdateSecurityGroupRequest) (*SecurityGroup, error) {
	var securityGroup *SecurityGroup
	req := c.R().
		SetBody(update).SetResult(&securityGroup)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", SecurityGroupEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return securityGroup, err
	}
	return securityGroup, nil
}

// DeleteSecurityGroup deletes a specific security group by its identity.
func (c *Client) DeleteSecurityGroup(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", SecurityGroupEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
