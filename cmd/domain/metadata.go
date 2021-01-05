package domain

import "time"

// Metadata represents the metadata of a secure domain.
type Metadata struct {
	Name                  string
	CertificateValidUntil time.Time
}
