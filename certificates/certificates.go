package certificates

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListCertificates lists certificates in a region.
func (c *Client) ListCertificates(ctx context.Context, region string) ([]Certificate, error) {
	items := []Certificate{}
	resp, err := c.Do(ctx, c.R().SetResult(&items), client.GET, regionPath(region, "certificates"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return items, err
	}
	return items, nil
}

// GetCertificate returns one certificate. The private key is not part of the response.
func (c *Client) GetCertificate(ctx context.Context, region, identity string) (*Certificate, error) {
	var item Certificate
	resp, err := c.Do(ctx, c.R().SetResult(&item), client.GET, regionPath(region, "certificates", identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// CreateCertificate starts issuance. Status is issuing until the order completes.
func (c *Client) CreateCertificate(ctx context.Context, region string, create CreateCertificateRequest) (*Certificate, error) {
	var item Certificate
	resp, err := c.Do(ctx, c.R().SetBody(create).SetResult(&item), client.POST, regionPath(region, "certificates"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// DeleteCertificate soft-deletes a certificate and its managed key.
func (c *Client) DeleteCertificate(ctx context.Context, region, identity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, regionPath(region, "certificates", identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// RotateCertificate starts a re-issue. The current PEM is kept until the new order succeeds.
func (c *Client) RotateCertificate(ctx context.Context, region, identity string) (*Certificate, error) {
	var item Certificate
	resp, err := c.Do(ctx, c.R().SetResult(&item), client.POST, regionPath(region, "certificates", identity, "rotate"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// RevokeCertificate marks a certificate revoked.
func (c *Client) RevokeCertificate(ctx context.Context, region, identity string) (*Certificate, error) {
	var item Certificate
	resp, err := c.Do(ctx, c.R().SetResult(&item), client.POST, regionPath(region, "certificates", identity, "revoke"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// ExportPrivateKey downloads the private key. The certificate must have allowPrivateKeyExport set.
func (c *Client) ExportPrivateKey(ctx context.Context, region, identity string) (*PrivateKey, error) {
	var body PrivateKey
	resp, err := c.Do(ctx, c.R().SetResult(&body), client.POST, regionPath(region, "certificates", identity, "private-key"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &body, err
	}
	decoded, err := decodeMaterial(body.PrivateKey)
	if err != nil {
		return nil, err
	}
	body.PrivateKey = decoded
	return &body, nil
}

// CertificatePEM downloads the public certificate. The response has no private key.
func (c *Client) CertificatePEM(ctx context.Context, region, identity string) (*PEM, error) {
	var body PEM
	resp, err := c.Do(ctx, c.R().SetResult(&body), client.GET, regionPath(region, "certificates", identity, "pem"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &body, err
	}
	return decodePEM(&body)
}

// ListAuthorities lists issuers in a region.
func (c *Client) ListAuthorities(ctx context.Context, region string) ([]CertificateAuthority, error) {
	items := []CertificateAuthority{}
	resp, err := c.Do(ctx, c.R().SetResult(&items), client.GET, regionPath(region, "issuers"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return items, err
	}
	return items, nil
}

// GetAuthority returns one issuer.
func (c *Client) GetAuthority(ctx context.Context, region, identity string) (*CertificateAuthority, error) {
	var item CertificateAuthority
	resp, err := c.Do(ctx, c.R().SetResult(&item), client.GET, regionPath(region, "issuers", identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// CreateAuthority provisions an issuer.
func (c *Client) CreateAuthority(ctx context.Context, region string, create CreateAuthorityRequest) (*CertificateAuthority, error) {
	var item CertificateAuthority
	resp, err := c.Do(ctx, c.R().SetBody(create).SetResult(&item), client.POST, regionPath(region, "issuers"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// DeleteAuthority deletes an issuer that has no leaves or cross-signed certificates.
func (c *Client) DeleteAuthority(ctx context.Context, region, identity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, regionPath(region, "issuers", identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// AuthorityPEM downloads the issuer certificate.
func (c *Client) AuthorityPEM(ctx context.Context, region, identity string) (*AuthorityPEM, error) {
	var body AuthorityPEM
	resp, err := c.Do(ctx, c.R().SetResult(&body), client.GET, regionPath(region, "issuers", identity, "pem"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &body, err
	}
	decoded, err := decodeMaterial(body.Certificate)
	if err != nil {
		return nil, err
	}
	body.Certificate = decoded
	return &body, nil
}

// RotateAuthority issues a new issuer key and keeps the previous certificate.
func (c *Client) RotateAuthority(ctx context.Context, region, identity string, rotate RotateAuthorityRequest) (*CertificateAuthority, error) {
	var item CertificateAuthority
	resp, err := c.Do(ctx, c.R().SetBody(rotate).SetResult(&item), client.POST, regionPath(region, "issuers", identity, "rotate"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// CrossSignAuthority stores an extra certificate for the current issuer public key, signed by another root.
func (c *Client) CrossSignAuthority(ctx context.Context, region, identity string, crossSign CrossSignAuthorityRequest) (*CertificateAuthority, error) {
	var item CertificateAuthority
	resp, err := c.Do(ctx, c.R().SetBody(crossSign).SetResult(&item), client.POST, regionPath(region, "issuers", identity, "cross-sign"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &item, err
	}
	return &item, nil
}

// AuthorityCertificatePEM downloads one historical or cross-signed issuer certificate.
func (c *Client) AuthorityCertificatePEM(ctx context.Context, region, identity, certificateIdentity string) (*AuthorityPEM, error) {
	var body AuthorityPEM
	resp, err := c.Do(ctx, c.R().SetResult(&body), client.GET, regionPath(region, "issuers", identity, "certificates", certificateIdentity, "pem"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &body, err
	}
	decoded, err := decodeMaterial(body.Certificate)
	if err != nil {
		return nil, err
	}
	body.Certificate = decoded
	return &body, nil
}

func decodePEM(body *PEM) (*PEM, error) {
	certificate, err := decodeMaterial(body.Certificate)
	if err != nil {
		return nil, err
	}
	chain, err := decodeMaterial(body.Chain)
	if err != nil {
		return nil, err
	}
	body.Certificate = certificate
	body.Chain = chain
	return body, nil
}

func decodeMaterial(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode certificate material: %w", err)
	}
	return string(raw), nil
}
