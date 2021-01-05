// Package domain provides a representation of a domain, which TLS certificate status can be checked.
package domain

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/gonzalo-bulnes/scanner"
)

// Domain represents a secure domain.
type Domain struct {
	getCertificate func(string) (*x509.Certificate, error)
	Name           string
}

// New returns a new domain.
func New(name string) *Domain {
	return &Domain{
		Name: name,
	}
}

// GetCertificateExpirationTime fetches the expiration time of the domain's TLS certificate.
func (d *Domain) GetCertificateExpirationTime(ctx context.Context) (validUntil time.Time, err error) {

	url := fmt.Sprintf("%s:443", d.Name)
	conn, err := tls.Dial("tcp", url, nil)
	if err != nil {
		err = fmt.Errorf("connection error %w", err)
	}

	// The first element is the leaf certificate
	// that the connection is verified against.
	//
	// See https://golang.org/pkg/crypto/tls
	validUntil = conn.ConnectionState().PeerCertificates[0].NotAfter
	return
}

// Check implements the scanner.Service interface.
func (d *Domain) Check(ctx context.Context) scanner.Status {

	validUntil, err := d.GetCertificateExpirationTime(ctx)
	if err != nil {
		return Status{err: fmt.Errorf("status check error: %w", err)}
	}

	return Status{
		metadata: Metadata{
			Name: d.Name,
			CertificateValidUntil: validUntil,
		},
	}
}
