package iaas

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

type Region struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	ObjectVersion int               `json:"objectVersion"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	Zones         []Zone            `json:"zones"`
}

type Zone struct {
	Identity            string            `json:"identity"`
	Name                string            `json:"name"`
	Slug                string            `json:"slug"`
	Description         string            `json:"description"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
	ObjectVersion       int               `json:"objectVersion"`
	Labels              map[string]string `json:"labels"`
	Annotations         map[string]string `json:"annotations"`
	CloudRegionIdentity string            `json:"cloudRegionIdentity"`
	CloudRegion         *Region           `json:"CloudRegion"`
}

// Stub for VpcFirewallRule (not given in your code)
type VpcFirewallRule struct {
	// Add fields if you have them defined elsewhere
}

type Vpc struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`
	Status        string    `json:"status"`

	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CIDRs       []string          `json:"cidrs"`

	Organisation  *base.Organisation `json:"organisation"`
	CloudRegion   *Region            `json:"cloudRegion"`
	Subnets       []Subnet           `json:"subnets"`
	FirewallRules []VpcFirewallRule  `json:"firewallRules"`
}

// Subnet
type Subnet struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`

	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`

	VpcIdentity string `json:"vpcIdentity"`
	Vpc         *Vpc   `json:"vpc"`
	CloudZone   *Zone  `json:"cloudZone"`
	Cidr        string `json:"cidr"`

	RouteTable     *RouteTable `json:"routeTable,omitempty"`
	V4usingIPs     int         `json:"v4usingIPs"`
	V4availableIPs int         `json:"v4availableIPs"`
	V6usingIPs     int         `json:"v6usingIPs"`
	V6availableIPs int         `json:"v6availableIPs"`
}

type VpcNatGateway struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`
	Status        string    `json:"status"`

	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`

	Organisation   *base.Organisation `json:"organisation"`
	VpcIdentity    string             `json:"vpcIdentity"`
	Vpc            *Vpc               `json:"vpc"`
	SubnetIdentity string             `json:"subnetIdentity"`
	Subnet         *Subnet            `json:"subnet"`
	EndpointIP     string             `json:"endpointIP"`

	V4IP string `json:"v4IP"`
	V6IP string `json:"v6IP"`
}

type VpcLoadbalancer struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	ObjectVersion int               `json:"objectVersion"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	Status        string            `json:"status"`

	Organisation   *base.Organisation `json:"organisation"`
	VpcIdentity    string             `json:"vpcIdentity"`
	Vpc            *Vpc               `json:"vpc"`
	SubnetIdentity string             `json:"subnetIdentity"`
	Subnet         *Subnet            `json:"subnet"`

	ExternalIpAddresses []string `json:"externalIpAddresses"`
	InternalIpAddresses []string `json:"internalIpAddresses"`
	Hostname            string   `json:"hostname"`

	LoadbalancerListeners []VpcLoadbalancerListener `json:"loadbalancerListeners"`
}

type VpcLoadbalancerListener struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`

	Port           int                         `json:"port"`
	Protocol       string                      `json:"protocol"`
	TargetGroup    *VpcLoadbalancerTargetGroup `json:"targetGroup"`
	TargetGroupId  int                         `json:"targetGroupId"`
	AllowedSources []string                    `json:"allowedSources"`
}

type VpcLoadbalancerTargetGroup struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`

	Organisation   *base.Organisation `json:"organisation"`
	Vpc            *Vpc               `json:"vpc"`
	TargetPort     int                `json:"targetPort"`
	Protocol       string             `json:"protocol"`
	TargetSelector map[string]string  `json:"targetSelector"`

	LoadbalancerListeners              []VpcLoadbalancerListener           `json:"loadbalancerListeners"`
	LoadbalancerTargetGroupAttachments []LoadbalancerTargetGroupAttachment `json:"loadbalancerTargetGroupAttachments"`
}

type LoadbalancerTargetGroupAttachment struct {
	Identity                  string                      `json:"identity"`
	CreatedAt                 time.Time                   `json:"createdAt"`
	LoadbalancerTargetGroupId int                         `json:"loadbalancerTargetGroupId"`
	LoadbalancerTargetGroup   *VpcLoadbalancerTargetGroup `json:"loadbalancerTargetGroup"`
	VirtualMachineInstanceId  int                         `json:"virtualMachineInstanceId"`
	VirtualMachineInstance    *Machine                    `json:"virtualMachineInstance"`
}

type Volume struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`
	Status        string    `json:"status"`

	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`

	// SourceMachineImage is the machine image that was used to create the volume. Only set if the volume was created from a machine image.
	SourceMachineImage *MachineImage      `json:"sourceMachineImage"`
	VolumeType         *VolumeType        `json:"volumeType"`
	Attachments        []VolumeAttachment `json:"attachments"`
	Organisation       *base.Organisation `json:"organisation"`
	CloudRegion        *Region            `json:"cloudRegion"`
	Size               int                `json:"size"`
	DeleteProtection   bool               `json:"deleteProtection"`
}

type VolumeType struct {
	Identity    string `json:"identity"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StorageType string `json:"storageType"`
	AllowResize bool   `json:"allowResize"`
}

type VolumeAttachment struct {
	Identity               string     `json:"identity"`
	CreatedAt              time.Time  `json:"createdAt"`
	Description            string     `json:"description"`
	DeviceName             string     `json:"deviceName"`
	AttachedToIdentity     string     `json:"attachedToIdentity"`
	AttachedToResourceType string     `json:"attachedToResourceType"`
	DetachmentRequestedAt  *time.Time `json:"detachmentRequestedAt,omitempty"`
	CanDetach              bool       `json:"canDetach"`

	// Only set if attachedToResourceType == "cloud_virtual_machine"
	VirtualMachine *Machine `json:"virtualMachine"`

	PersistentVolume *Volume `json:"persistentVolume"`
}

type VpcGatewayEndpoint struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   *string           `json:"description,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`           // ISO date string
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"` // optional
	DeletedAt     *time.Time        `json:"deletedAt,omitempty"` // optional
	ObjectVersion int               `json:"objectVersion"`

	EndpointAddress  string             `json:"endpointAddress"` // IP address
	EndpointHostname string             `json:"endpointHostname"`
	Vpc              *Vpc               `json:"vpc,omitempty"`
	Organisation     *base.Organisation `json:"organisation,omitempty"`
	CloudRegion      *Region            `json:"cloudRegion,omitempty"`
	Subnet           *Subnet            `json:"subnet,omitempty"`
	Status           string             `json:"status"`
}

type RouteTable struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   *string           `json:"description,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`           // ISO date string
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"` // optional
	DeletedAt     *time.Time        `json:"deletedAt,omitempty"` // optional
	ObjectVersion int               `json:"objectVersion"`

	Organisation      *base.Organisation `json:"organisation,omitempty"`
	Vpc               *Vpc               `json:"vpc"`
	Routes            []RouteEntry       `json:"routes,omitempty"`
	IsDefault         bool               `json:"isDefault"`
	AssociatedSubnets []Subnet           `json:"associatedSubnets"`
}

type RouteEntry struct {
	Identity                 string              `json:"identity"`
	Note                     *string             `json:"note,omitempty"`
	RouteTable               *RouteTable         `json:"routeTable,omitempty"`
	DestinationCidrBlock     string              `json:"destinationCidrBlock"`
	TargetGatewayIdentity    *string             `json:"targetGatewayIdentity,omitempty"`
	TargetGateway            *VpcGatewayEndpoint `json:"targetGateway,omitempty"`
	TargetNatGatewayIdentity *string             `json:"targetNatGatewayIdentity,omitempty"`
	TargetNatGateway         *VpcNatGateway      `json:"targetNatGateway,omitempty"`
	GatewayAddress           *string             `json:"gatewayAddress,omitempty"`
	TargetGatewayEndpoint    *VpcGatewayEndpoint `json:"targetGatewayEndpoint,omitempty"`
	Type                     string              `json:"type"`
}

type Machine struct {
	Identity         string            `json:"identity"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        *time.Time        `json:"updatedAt,omitempty"`
	Description      *string           `json:"description,omitempty"`
	Annotations      map[string]string `json:"annotations,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	State            MachineState
	CloudInit        *string `json:"cloudInit"`
	DeleteProtection bool    `json:"deleteProtection"`
	// SecurityGroups    []SecurityGroup          `json:"securityGroups,omitempty"`
	Organisation      *base.Organisation       `json:"organisation,omitempty"`
	MachineType       *MachineType             `json:"machineType,omitempty"`
	MachineImage      *MachineImage            `json:"machineImage,omitempty"`
	PersistentVolume  *Volume                  `json:"persistentVolume,omitempty" validate:"-"`
	Vpc               *Vpc                     `json:"vpc,omitempty"`
	Subnet            *Subnet                  `json:"subnet,omitempty"`
	Interfaces        VirtualMachineInterfaces `json:"interfaces,omitempty"`
	VolumeAttachments []VolumeAttachment       `json:"volumeAttachments,omitempty"`
	Status            ResourceStatus           `json:"status"`
}

type ResourceStatus struct {
	Status             string    `json:"status"`
	StatusMessage      string    `json:"statusMessage"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
}

type VirtualMachineInterfaces []VirtualMachineInterface
type VirtualMachineInterface struct {
	// Name is the name of the interface
	Name string `json:"name" validate:"required"`
	// MacAddress is the MAC address of the interface
	MacAddress string `json:"macAddress"`
	// IPAddresses is a list of IP addresses that are assigned to the interface
	IPAddresses []string `json:"ipAddresses"`
}

type MachineState string

const (
	// MachineStateCreating is the state of the machine that is being created
	MachineStateCreating MachineState = "creating"
	// MachineStateRunning is the state of the machine that is running
	MachineStateRunning MachineState = "running"
	// MachineStateStopped is the state of the machine that is stopped
	MachineStateStopped MachineState = "stopped"
	// MachineStateDeleting is the state of the machine that is being deleted
	MachineStateDeleting MachineState = "deleting"
	// MachineStateDeleted is the state of the machine that is deleted
	MachineStateDeleted MachineState = "deleted"
)

type MachineTypeCategory struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	MachineTypes []MachineType `json:"machineTypes"`
}

type MachineType struct {
	Identity    string `json:"identity"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Vcpus       int    `json:"vcpus"`
	RamMb       int    `json:"ramMb"`
	DiskGb      int    `json:"diskGb"`
	SwapMb      int    `json:"swapMb"`
}

type MachineImage struct {
	Identity     string            `json:"identity"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	Labels       map[string]string `json:"labels"`
	Description  string            `json:"description"`
	Architecture string            `json:"architecture"`
}

type CreateVpc struct {
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	Labels              map[string]string `json:"labels"`
	Annotations         map[string]string `json:"annotations"`
	CloudRegionIdentity string            `json:"cloudRegionIdentity"`
	VpcCidrs            []string          `json:"vpcCidrs"`
}

type UpdateVpc struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	VpcCidrs    []string          `json:"vpcCidrs"`
}

type CreateVolume struct {
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	Labels              map[string]string `json:"labels"`
	Annotations         map[string]string `json:"annotations"`
	Type                string            `json:"type"`
	Size                int               `json:"size"`
	CloudRegionIdentity string            `json:"cloudRegionIdentity"`
	VolumeTypeIdentity  string            `json:"volumeTypeIdentity"`
	DeleteProtection    bool              `json:"deleteProtection"`
}

type UpdateVolume struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Labels           map[string]string `json:"labels"`
	Annotations      map[string]string `json:"annotations"`
	Size             int               `json:"size"`
	DeleteProtection bool              `json:"deleteProtection"`
}

type AttachVolumeRequest struct {
	Description      string `json:"description"`
	DeviceName       string `json:"deviceName"`
	ResourceType     string `json:"resourceType"`
	ResourceIdentity string `json:"resourceIdentity"`
}

type DetachVolumeRequest struct {
	ResourceType     string `json:"resourceType"`
	ResourceIdentity string `json:"resourceIdentity"`
}

type CreateSubnet struct {
	Name                         string            `json:"name"`
	Description                  string            `json:"description"`
	Labels                       map[string]string `json:"labels,omitempty"`
	Annotations                  map[string]string `json:"annotations,omitempty"`
	VpcIdentity                  string            `json:"vpcIdentity"`
	CloudZone                    string            `json:"cloudZone"`
	Cidr                         string            `json:"cidr"`
	AssociatedRouteTableIdentity *string           `json:"associatedRouteTableIdentity,omitempty"`
}

type UpdateSubnet struct {
	Name                         string            `json:"name"`
	Description                  string            `json:"description"`
	Labels                       map[string]string `json:"labels,omitempty"`
	Annotations                  map[string]string `json:"annotations,omitempty"`
	AssociatedRouteTableIdentity *string           `json:"associatedRouteTableIdentity,omitempty"`
}

type UpdateRouteTableRoutes struct {
	Routes []UpdateRouteTableRoute `json:"routes"`
}

type CreateRouteTableRoute struct {
	DestinationCidrBlock     string `json:"destinationCidrBlock"`
	TargetGatewayIdentity    string `json:"targetGatewayIdentity,omitempty"`
	TargetNatGatewayIdentity string `json:"targetNatGatewayIdentity,omitempty"`
	GatewayAddress           string `json:"gatewayAddress,omitempty"`
}

type UpdateRouteTableRoute struct {
	DestinationCidrBlock     string `json:"destinationCidrBlock"`
	TargetGatewayIdentity    string `json:"targetGatewayIdentity,omitempty"`
	TargetNatGatewayIdentity string `json:"targetNatGatewayIdentity,omitempty"`
	GatewayAddress           string `json:"gatewayAddress,omitempty"`
}

type CreateRouteTable struct {
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	VpcIdentity string            `json:"vpcIdentity"`
}

type UpdateRouteTable struct {
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type CreateVpcNatGateway struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Labels         map[string]string `json:"labels"`
	Annotations    map[string]string `json:"annotations"`
	VpcIdentity    string            `json:"vpcIdentity"`
	SubnetIdentity string            `json:"subnetIdentity"`
}

type UpdateVpcNatGateway struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type CreateMachine struct {
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	Labels           map[string]string   `json:"labels"`
	Annotations      map[string]string   `json:"annotations"`
	Subnet           string              `json:"subnet"`
	CloudInit        string              `json:"cloudInit"`
	CloudInitRef     string              `json:"cloudInitRef"`
	MachineImage     string              `json:"machineImage"`
	MachineType      string              `json:"machineType"`
	DeleteProtection bool                `json:"deleteProtection"`
	VpcIdentity      string              `json:"vpcIdentity"`
	RootVolume       CreateMachineVolume `json:"rootVolume"`
}

type CreateMachineVolume struct {
	ExistingVolumeRef  string `json:"existingVolumeRef"`
	VolumeTypeIdentity string `json:"volumeTypeIdentity"`
	Size               int    `json:"size"`
	Name               string `json:"name"`
	Description        string `json:"description"`
}

type UpdateMachine struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Labels           map[string]string `json:"labels"`
	Annotations      map[string]string `json:"annotations"`
	Subnet           string            `json:"subnet"`
	MachineType      string            `json:"machineType"`
	DeleteProtection bool              `json:"deleteProtection"`
}

type CreateLoadbalancer struct {
	Name                     string            `json:"name"`
	Description              string            `json:"description"`
	Labels                   map[string]string `json:"labels,omitempty"`
	Annotations              map[string]string `json:"annotations,omitempty"`
	Subnet                   string            `json:"subnet"`
	InternalLoadbalancer     bool              `json:"internalLoadbalancer"`
	DeleteProtection         bool              `json:"deleteProtection"`
	Listeners                []CreateListener  `json:"listeners"`
	SecurityGroupAttachments []string          `json:"securityGroupAttachments,omitempty"`
}

type UpdateLoadbalancer struct {
	Name                     string            `json:"name"`
	Description              string            `json:"description"`
	Labels                   map[string]string `json:"labels,omitempty"`
	Annotations              map[string]string `json:"annotations,omitempty"`
	DeleteProtection         bool              `json:"deleteProtection"`
	SecurityGroupAttachments []string          `json:"securityGroupAttachments,omitempty"`
}

type CreateTargetGroup struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	Vpc            string            `json:"vpc"`
	TargetPort     int               `json:"targetPort"`
	Protocol       string            `json:"protocol"`
	TargetSelector map[string]string `json:"targetSelector,omitempty"`
}

type UpdateTargetGroup struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	TargetPort     int               `json:"targetPort"`
	Protocol       string            `json:"protocol"`
	TargetSelector map[string]string `json:"targetSelector,omitempty"`
}

type AttachTargetRequest struct {
	ServerIdentity string `json:"serverIdentity"`
}

type CreateListener struct {
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	Port           int               `json:"port"`
	Protocol       string            `json:"protocol"`
	TargetGroup    string            `json:"targetGroup"`
	AllowedSources []string          `json:"allowedSources,omitempty"`
}

type UpdateListener struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Port        int               `json:"port"`
	Protocol    string            `json:"protocol"`
	TargetGroup string            `json:"targetGroup"`
}
