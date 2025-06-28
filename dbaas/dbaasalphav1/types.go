package dbaasalphav1

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type DbClusterDatabaseEngine string

const (
	DbClusterDatabaseEnginePostgres DbClusterDatabaseEngine = "postgres"
)

type DbCluster struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	ObjectVersion int         `json:"objectVersion"`
	Labels        Labels      `json:"labels"`
	Annotations   Annotations `json:"annotations"`

	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Vpc is the VPC the cluster is deployed in
	Vpc    *iaas.Vpc    `json:"vpc,omitempty"`
	Region *iaas.Region `json:"region,omitempty"`
	// Subnet is the subnet the cluster is deployed in
	Subnet *iaas.Subnet `json:"subnet,omitempty"`
	// DatabaseInstanceType is the instance type used to determine the size of the cluster instances
	DatabaseInstanceType *iaas.MachineType `json:"database_instance_type,omitempty"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
	// Engine is the database engine of the cluster
	Engine DbClusterDatabaseEngine `json:"engine"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// DatabaseEngineVersion is the version of the database engine
	DatabaseEngineVersion *DbClusterEngineVersion `json:"database_engine_version,omitempty"`
	// // DbParameterGroupId is the ID of the database parameter group
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// VolumeTypeClass is the storage type used to determine the size of the cluster storage
	VolumeTypeClass *iaas.VolumeType `json:"volume_type_class,omitempty"`
	// AutoMinorVersionUpgrade is a flag indicating if the cluster should automatically upgrade to the latest minor version
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`
	// DatabaseName is the name of the database on the cluster. Optional name. If provided, it will be used as the name of the database on the cluster.
	DatabaseName *string `json:"databaseName"`
	// DeleteProtection is a flag indicating if the cluster should be protected from deletion. The database cannot be deleted if this is true.
	DeleteProtection bool `json:"deleteProtection"`
	// SecurityGroups is a list of security groups associated with the cluster
	SecurityGroups []iaas.SecurityGroup `json:"securityGroups,omitempty"`
	// Status is the status of the cluster
	Status DbClusterStatus `json:"status"`
	// EndpointIpv4 is the IPv4 address of the cluster endpoint
	EndpointIpv4 string `json:"endpointIpv4"`
	// EndpointIpv6 is the IPv6 address of the cluster endpoint
	EndpointIpv6 string `json:"endpointIpv6"`
	// Port is the port of the cluster endpoint
	Port int `json:"port"`
}

type DbClusterEngineVersion struct {
	Identity string `json:"identity"`
	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`
	// Engine is the database engine
	Engine DbClusterDatabaseEngine `json:"engine"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// MajorVersion is the major version of the engine
	MajorVersion int `json:"majorVersion"`
	// MinorVersion is the minor version of the engine
	MinorVersion int `json:"minorVersion"`
	// Supported is a flag indicating if the engine version is supported
	Supported bool `json:"supported"`
	// MinMajorVersionUpgradeFrom is the minimum major version required to upgrade from
	MinMajorVersionUpgradeFrom *int `json:"minMajorVersionUpgradeFrom"`
	// MinMinorVersionUpgradeFrom is the minimum minor version required to upgrade from
	MinMinorVersionUpgradeFrom *int `json:"minMinorVersionUpgradeFrom"`
	// MaxMajorVersionUpgradeTo is the maximum major version that can be upgraded to
	MaxMajorVersionUpgradeTo *int `json:"maxMajorVersionUpgradeTo"`
	// MaxMinorVersionUpgradeTo is the maximum minor version that can be upgraded to
	MaxMinorVersionUpgradeTo *int `json:"maxMinorVersionUpgradeTo"`
	// Enabled is a flag indicating if the engine version is enabled
	Enabled bool `json:"enabled"`
	// DefaultParameters is a map of parameter name to database engine specific parameter value
	DefaultParameters map[string]string `json:"defaultParameters"`
}

type CreateDbClusterRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	// Subnet is the subnet identity of the cloud subnet
	SubnetIdentity           string   `json:"subnetIdentity"`
	SecurityGroupAttachments []string `json:"securityGroupAttachments"`
	DeleteProtection         bool     `json:"deleteProtection"`
	// Engine is the database engine
	Engine DbClusterDatabaseEngine `json:"engine"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// VolumeTypeClassIdentity is the identity of the storage type
	VolumeTypeClassIdentity string `json:"volumeTypeClassIdentity"`
	// DatabaseInstanceTypeIdentity is the identity of the database instance type
	DatabaseInstanceTypeIdentity string `json:"databaseInstanceTypeIdentity"`
	// AutoMinorVersionUpgrade is a flag indicating if the cluster should automatically upgrade to the latest minor version
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`
	// Instances is the number of instances in the cluster
	Instances int `json:"replicas"`
	// PostgresInitDb is the initial database to create on the cluster
	PostgresInitDb *PostgresInitDb `json:"postgresInitDb,omitempty"`

	// RestoreFromBackupIdentity is the identity of the backup to restore from
	RestoreFromBackupIdentity *string `json:"restoreFromBackupIdentity,omitempty"`
}

type PostgresInitDb struct {
	// DataChecksums is a flag to indicate if data checksums should be enabled
	DataChecksums bool `json:"dataChecksums,omitempty"`
	// Maps to the `ENCODING` parameter of `CREATE DATABASE`. This setting
	// cannot be changed. Character set encoding to use in the database.
	Encoding string `json:"encoding,omitempty"`

	// Maps to the `LOCALE` parameter of `CREATE DATABASE`. This setting
	// cannot be changed. Sets the default collation order and character
	// classification in the new database.
	Locale string `json:"locale,omitempty"`

	// Maps to the `LOCALE_PROVIDER` parameter of `CREATE DATABASE`. This
	// setting cannot be changed. This option sets the locale provider for
	// databases created in the new cluster. Available from PostgreSQL 16.
	LocaleProvider string `json:"localeProvider,omitempty"`

	// Maps to the `LC_COLLATE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	LcCollate string `json:"localeCollate,omitempty"`

	// Maps to the `LC_CTYPE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	LcCtype string `json:"localeCType,omitempty"`

	// Maps to the `ICU_LOCALE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// Specifies the ICU locale when the ICU provider is used.
	// This option requires `localeProvider` to be set to `icu`. Available from PostgreSQL 15.
	IcuLocale string `json:"icuLocale,omitempty"`

	// Maps to the `ICU_RULES` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// Specifies additional collation rules to customize the behavior of the default collation.
	// This option requires `localeProvider` to be set to `icu`. Available from PostgreSQL 16.
	IcuRules string `json:"icuRules,omitempty"`

	// Maps to the `BUILTIN_LOCALE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// Specifies the locale name when the builtin provider is used. This option requires `localeProvider` to be set to `builtin`.
	// Available from PostgreSQL 17.
	BuiltinLocale string `json:"builtinLocale,omitempty"`
	// Maps to the `COLLATION_VERSION` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// CollationVersion string `json:"collationVersion,omitempty"`
	// The value in megabytes (1 to 1024) to be passed to the `--wal-segsize`
	// option for initdb (default: empty, resulting in PostgreSQL default: 16MB)
	// +optional
	WalSegmentSize int `json:"walSegmentSize,omitempty"`
}

type UpdateDbClusterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels                   Labels   `json:"labels,omitempty"`
	SecurityGroupAttachments []string `json:"securityGroupAttachments"`
	DeleteProtection         bool     `json:"deleteProtection"`
	// EngineVersion is the version of the database engine
	EngineVersion *string `json:"engineVersion,omitempty"`
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// AutoMinorVersionUpgrade is a flag indicating if the cluster should automatically upgrade to the latest minor version
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`
	// DatabaseName is the name of the database on the cluster. Optional name. If provided, it will be used as the name of the database on the cluster.
	DatabaseName *string `json:"databaseName"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
	// DatabaseInstanceTypeIdentity is the identity of the database instance type. Optional identity. If provided, it will be used as the database instance type for the cluster.
	DatabaseInstanceTypeIdentity *string `json:"databaseInstanceTypeIdentity"`
}

type DbClusterStatus string

const (
	DbClusterStatusPending               DbClusterStatus = "pending"
	DbClusterStatusCreating              DbClusterStatus = "creating"
	DbClusterStatusReady                 DbClusterStatus = "ready"
	DbClusterStatusUpdating              DbClusterStatus = "updating"
	DbClusterStatusUpgradingMajorVersion DbClusterStatus = "upgrading-major-version"
	DbClusterStatusUpgradingMinorVersion DbClusterStatus = "upgrading-minor-version"
	DbClusterStatusFailed                DbClusterStatus = "failed"
	DbClusterStatusDeleting              DbClusterStatus = "deleting"
	DbClusterStatusDeleted               DbClusterStatus = "deleted"
	DbClusterStatusUnknown               DbClusterStatus = "unknown"
)

type CreatePgDatabaseRequest struct {
	// Name is the name of the database.
	Name string `json:"name"`
	// Maps to the `OWNER` parameter of `CREATE DATABASE`.
	// Maps to the `OWNER TO` command of `ALTER DATABASE`.
	// The role name of the user who owns the database inside PostgreSQL.
	Owner string `json:"owner"`

	// Maps to the `ALLOW_CONNECTIONS` parameter of `CREATE DATABASE` and
	// `ALTER DATABASE`. If false then no one can connect to this database.
	AllowConnections *bool `json:"allowConnections,omitempty"`

	// Maps to the `CONNECTION LIMIT` clause of `CREATE DATABASE` and
	// `ALTER DATABASE`. How many concurrent connections can be made to
	// this database. -1 (the default) means no limit.
	ConnectionLimit *int `json:"connectionLimit,omitempty"`
}

type UpdatePgDatabaseRequest struct {
	// Maps to the `ALLOW_CONNECTIONS` parameter of `CREATE DATABASE` and
	// `ALTER DATABASE`. If false then no one can connect to this database.
	AllowConnections *bool `json:"allowConnections,omitempty"`

	// Maps to the `CONNECTION LIMIT` clause of `CREATE DATABASE` and
	// `ALTER DATABASE`. How many concurrent connections can be made to
	// this database. -1 (the default) means no limit.
	ConnectionLimit *int `json:"connectionLimit,omitempty"`
}

type CreatePgRoleRequest struct {
	Name string `json:"name"`
	// Login is a flag to indicate if the role can login
	Login bool `json:"login"`

	// CreateDb is a flag to indicate if the role can create databases
	CreateDb bool `json:"createDb"`

	// CreateRole is a flag to indicate if the role can create roles
	CreateRole bool `json:"createRole"`

	// ConnectionLimit is the maximum number of concurrent connections for the role. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`

	// ValidUntil is the date and time the role will expire
	ValidUntil *time.Time `json:"validUntil,omitempty"`

	// Password is the password for the role
	Password string `json:"password,omitempty"`
}

type UpdatePgRoleRequest struct {
	// ConnectionLimit is the maximum number of concurrent connections for the role. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`

	// ValidUntil is the date and time the role will expire
	ValidUntil *time.Time `json:"validUntil,omitempty"`

	// Password is the password for the role. If provided, the password will be updated. If not provided, the password will not be updated.
	Password *string `json:"password,omitempty"`
}

type CreatePgBackupScheduleRequest struct {
	Name            string      `json:"name"`
	Description     *string     `json:"description,omitempty"`
	Annotations     Annotations `json:"annotations,omitempty"`
	Labels          Labels      `json:"labels,omitempty"`
	Schedule        string      `json:"schedule"`
	RetentionPolicy string      `json:"retentionPolicy"`
	Target          string      `json:"target,omitempty"`
}

type UpdatePgBackupScheduleRequest struct {
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Annotations     Annotations `json:"annotations,omitempty"`
	Labels          Labels      `json:"labels,omitempty"`
	Schedule        string      `json:"schedule"`
	RetentionPolicy string      `json:"retentionPolicy"`
	Target          string      `json:"target,omitempty"`
}

type DbClusterBackupScheduleMethod string

const (
	DbClusterBackupScheduleMethodSnapshot DbClusterBackupScheduleMethod = "snapshot"
)

type DbClusterBackupScheduleTarget string

const (
	DbClusterBackupScheduleTargetPrimary       DbClusterBackupScheduleTarget = "primary"
	DbClusterBackupScheduleTargetPreferStandby DbClusterBackupScheduleTarget = "prefer-standby"
)

type DbClusterBackupSchedule struct {
	// Identity is a unique identifier for the backup schedule
	Identity string `json:"identity"`

	// Name is a human-readable name of the backup schedule
	Name string `json:"name"`

	// Status is the status of the role
	Status ObjectStatus `json:"status"`

	// StatusMessage is the message of the role status
	StatusMessage string `json:"statusMessage,omitempty"`

	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`

	// DbCluster is the cluster the backup schedule belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`

	// Organisation is the organisation the backup schedule belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`

	// Method is the method of the backup schedule
	Method DbClusterBackupScheduleMethod `json:"method"`

	// Schedule is the schedule of the backup. Cron expression.
	// see https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
	Schedule string `json:"schedule"`

	// RetentionPolicy is the retention policy of the backup
	RetentionPolicy string `json:"retentionPolicy"`

	// NextBackupAt is the date and time the next backup will be taken
	NextBackupAt *time.Time `json:"nextBackupAt,omitempty"`

	// LastBackupAt is the date and time the last backup was taken
	LastBackupAt *time.Time `json:"lastBackupAt,omitempty"`

	// BackupCount is the number of backups
	BackupCount int64 `json:"backupCount"`

	// Suspended is a flag to indicate if the backup schedule is suspended
	Suspended bool `json:"suspended"`

	// Target is the target of the backup schedule
	Target DbClusterBackupScheduleTarget `json:"target"`

	// DeleteScheduledAt is the date and time the backup schedule will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}

type CreateDbClusterBackupRequest struct {
	Name            string      `json:"name"`
	Description     *string     `json:"description,omitempty"`
	Labels          Labels      `json:"labels"`
	Annotations     Annotations `json:"annotations"`
	RetentionPolicy *string     `json:"retentionPolicy,omitempty"`
	Target          string      `json:"target"`
}

type DbClusterBackup struct {
	// Identity is a unique identifier for the backup
	Identity string `json:"identity" validate:"required,notblank"`
	// DbCluster is the cluster the backup belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`
	// Organisation is the organisation the backup belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Labels is a map of labels for the backup
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is a map of annotations for the backup
	Annotations map[string]string `json:"annotations,omitempty"`
	// Region is the region the backup belongs to
	Region *iaas.Region `json:"region,omitempty"`
	// BackupSchedule is the backup schedule the backup belongs to
	// +optional
	BackupSchedule *DbClusterBackupSchedule `json:"backupSchedule,omitempty"`
	// BackupTrigger is the trigger of the backup
	BackupTrigger DbClusterBackupTrigger `json:"backupTrigger"`
	// EngineType is the type of the database engine, used for back-up restore purposes
	EngineType DbClusterDatabaseEngine `json:"engineType"`
	// EngineVersion is the version of the database engine, used for back-up restore purposes
	EngineVersion string `json:"engineVersion"`

	// BackupType is the type of the backup
	BackupType string `json:"backupType"`

	// Online is a flag to indicate if the backup is an online backup or offline/cold backup
	Online bool `json:"online"`

	// BeginLSN is the starting LSN of the backup
	BeginLSN string `json:"beginLSN"`

	// EndLSN is the ending LSN of the backup
	EndLSN string `json:"endLSN"`

	// BeginWAL is the starting WAL of the backup
	BeginWAL string `json:"beginWAL"`

	// EndWAL is the ending WAL of the backup
	EndWAL string `json:"endWAL"`

	// StartedAt is the date and time the backup started
	StartedAt *time.Time `json:"startedAt,omitempty"`

	// StoppedAt is the date and time the backup stopped
	StoppedAt *time.Time `json:"stoppedAt,omitempty"`

	// Status is the status of the backup
	Status ObjectStatus `json:"status"`

	// StatusMessage is the message of the backup status
	StatusMessage string `json:"statusMessage,omitempty"`

	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`

	// DeleteScheduledAt is the date and time the backup will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}

type DbClusterBackupTrigger string

const (
	DbClusterBackupTriggerManual   DbClusterBackupTrigger = "manual"
	DbClusterBackupTriggerSchedule DbClusterBackupTrigger = "schedule"
	DbClusterBackupTriggerSystem   DbClusterBackupTrigger = "system"
)

type CreatePgGrantRequest struct {
	Name         string `json:"name"`
	RoleName     string `json:"roleName"`
	DatabaseName string `json:"databaseName"`

	// Read is a flag to indicate if the role can read from the database
	Read bool `json:"read"`
	// Write is a flag to indicate if the role can write to the database
	Write bool `json:"write"`
}

type UpdatePgGrantRequest struct {
	Read  *bool `json:"read"`
	Write *bool `json:"write"`
}

type DbClusterPostgresRole struct {
	// Identity is a unique identifier for the role
	Identity string `json:"identity"`
	// Name is a human-readable name of the role
	Name string `json:"name"`
	// Status is the status of the role
	Status ObjectStatus `json:"status"`
	// StatusMessage is the message of the role status
	StatusMessage string `json:"statusMessage,omitempty"`
	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`
	// DbCluster is the cluster the role belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`
	// Login is a flag to indicate if the role can login
	Login bool `json:"login"`
	// CreateDb is a flag to indicate if the role can create databases
	CreateDb bool `json:"createDb"`
	// CreateRole is a flag to indicate if the role can create roles
	CreateRole bool `json:"createRole"`
	// ConnectionLimit is the maximum number of concurrent connections for the role. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`
	// ValidUntil is the date and time the role will expire
	ValidUntil *time.Time `json:"validUntil,omitempty"`
	// DeleteScheduledAt is the date and time the role will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}
