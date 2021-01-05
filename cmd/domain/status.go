package domain

// Status of a TLS certificate.
type Status struct {
	metadata Metadata
	err      error
}

// Value implements the scanner.Status interface.
func (s Status) Value() interface{} {
	return s.metadata
}

// Err implements the scanner.Status interface.
func (s Status) Err() error {
	return s.err
}
