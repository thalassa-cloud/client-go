package certificates

// Issuer is the issuer selected by a certificate.
type Issuer struct {
	Identity string `json:"identity,omitempty"`
	Name     string `json:"name,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// Subject is a distinguished name.
type Subject struct {
	CommonName         string `json:"commonName"`
	Organization       string `json:"organization,omitempty"`
	OrganizationalUnit string `json:"organizationalUnit,omitempty"`
	Country            string `json:"country,omitempty"`
}

// OtherName is a subjectAltName otherName.
type OtherName struct {
	OID   string `json:"oid"`
	Value string `json:"value"`
}

// Certificate is a TLS certificate. The private key is not included.
type Certificate struct {
	Identity              string      `json:"identity"`
	Name                  string      `json:"name"`
	Slug                  string      `json:"slug,omitempty"`
	Description           string      `json:"description,omitempty"`
	ProjectID             string      `json:"projectId,omitempty"`
	Status                string      `json:"status"`
	DNSNames              []string    `json:"dnsNames,omitempty"`
	Subject               *Subject    `json:"subject,omitempty"`
	IPAddresses           []string    `json:"ipAddresses,omitempty"`
	URIs                  []string    `json:"uris,omitempty"`
	OtherNames            []OtherName `json:"otherNames,omitempty"`
	KeyUsages             []string    `json:"keyUsages,omitempty"`
	ExtKeyUsages          []string    `json:"extKeyUsages,omitempty"`
	KeyType               string      `json:"keyType,omitempty"`
	SelfSigned            bool        `json:"selfSigned"`
	Issuer                Issuer      `json:"issuer,omitempty"`
	ZoneIdentity          string      `json:"zoneIdentity,omitempty"`
	SerialNumber          string      `json:"serialNumber,omitempty"`
	NotBefore             string      `json:"notBefore,omitempty"`
	NotAfter              string      `json:"notAfter,omitempty"`
	ValidityDays          int         `json:"validityDays,omitempty"`
	AutoRenew             bool        `json:"autoRenew"`
	AllowPrivateKeyExport bool        `json:"allowPrivateKeyExport"`
	RenewBeforeDays       int         `json:"renewBeforeDays,omitempty"`
	LastError             string      `json:"lastError,omitempty"`
	LastIssuedAt          string      `json:"lastIssuedAt,omitempty"`
	RevokedAt             string      `json:"revokedAt,omitempty"`
}

// CertificateAuthority is a customer issuer. The signing key is not included.
type CertificateAuthority struct {
	Identity       string                 `json:"identity"`
	Name           string                 `json:"name"`
	Slug           string                 `json:"slug,omitempty"`
	Description    string                 `json:"description,omitempty"`
	ProjectID      string                 `json:"projectId,omitempty"`
	Status         string                 `json:"status"`
	Kind           string                 `json:"kind,omitempty"`
	Provider       string                 `json:"provider,omitempty"`
	Subject        *Subject               `json:"subject,omitempty"`
	IPAddresses    []string               `json:"ipAddresses,omitempty"`
	URIs           []string               `json:"uris,omitempty"`
	OtherNames     []OtherName            `json:"otherNames,omitempty"`
	KeyUsages      []string               `json:"keyUsages,omitempty"`
	ExtKeyUsages   []string               `json:"extKeyUsages,omitempty"`
	KeyType        string                 `json:"keyType,omitempty"`
	Type           string                 `json:"type,omitempty"`
	ParentIdentity string                 `json:"parentIdentity,omitempty"`
	PathLen        int                    `json:"pathLen,omitempty"`
	SerialNumber   string                 `json:"serialNumber,omitempty"`
	NotBefore      string                 `json:"notBefore,omitempty"`
	NotAfter       string                 `json:"notAfter,omitempty"`
	ValidityDays   int                    `json:"validityDays,omitempty"`
	LastError      string                 `json:"lastError,omitempty"`
	Certificates   []AuthorityCertificate `json:"certificates,omitempty"`
}

// AuthorityCertificate is one public certificate for an issuer key.
type AuthorityCertificate struct {
	Identity       string `json:"identity"`
	IssuerIdentity string `json:"issuerIdentity,omitempty"`
	Kind           string `json:"kind"`
	Current        bool   `json:"current"`
	SerialNumber   string `json:"serialNumber,omitempty"`
	NotBefore      string `json:"notBefore,omitempty"`
	NotAfter       string `json:"notAfter,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
}

// PEM is a decoded certificate download. The HTTP body is standard base64. It has no private key.
type PEM struct {
	Certificate string `json:"certificate"`
	Chain       string `json:"chain,omitempty"`
}

// AuthorityPEM is an issuer certificate download.
type AuthorityPEM struct {
	Certificate string `json:"certificate"`
}

// CreateCertificateRequest starts an asynchronous issuance.
type CreateCertificateRequest struct {
	Name                  string      `json:"name"`
	Description           string      `json:"description,omitempty"`
	IssuerIdentity        string      `json:"issuerIdentity,omitempty"`
	SelfSigned            bool        `json:"selfSigned,omitempty"`
	DNSNames              []string    `json:"dnsNames,omitempty"`
	Subject               *Subject    `json:"subject,omitempty"`
	IPAddresses           []string    `json:"ipAddresses,omitempty"`
	URIs                  []string    `json:"uris,omitempty"`
	OtherNames            []OtherName `json:"otherNames,omitempty"`
	KeyUsages             []string    `json:"keyUsages,omitempty"`
	ExtKeyUsages          []string    `json:"extKeyUsages,omitempty"`
	ZoneIdentity          string      `json:"zoneIdentity,omitempty"`
	KeyType               string      `json:"keyType,omitempty"`
	ValidityDays          int         `json:"validityDays,omitempty"`
	AutoRenew             *bool       `json:"autoRenew,omitempty"`
	AllowPrivateKeyExport bool        `json:"allowPrivateKeyExport,omitempty"`
}

// PrivateKey is the decoded PKCS#8 PEM from ExportPrivateKey. The HTTP body is standard base64.
type PrivateKey struct {
	PrivateKey string `json:"privateKey"`
}

// CreateAuthorityRequest provisions an issuer.
type CreateAuthorityRequest struct {
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Kind         string      `json:"kind,omitempty"`
	Provider     string      `json:"provider,omitempty"`
	Subject      *Subject    `json:"subject,omitempty"`
	IPAddresses  []string    `json:"ipAddresses,omitempty"`
	URIs         []string    `json:"uris,omitempty"`
	OtherNames   []OtherName `json:"otherNames,omitempty"`
	KeyUsages    []string    `json:"keyUsages,omitempty"`
	ExtKeyUsages []string    `json:"extKeyUsages,omitempty"`
	KeyType      string      `json:"keyType,omitempty"`
	ValidityDays int         `json:"validityDays,omitempty"`
}

// RotateAuthorityRequest asks for a new issuer key. The previous certificate is kept.
type RotateAuthorityRequest struct{}

// CrossSignAuthorityRequest asks another root to sign this issuer's current public key.
type CrossSignAuthorityRequest struct {
	IssuerAuthorityIdentity string `json:"issuerIdentity"`
}
