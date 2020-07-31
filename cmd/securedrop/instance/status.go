package instance

// Status of a SecureDrop instance.
type Status struct {
	Metadata
	err error
}

// Value implements the scanner.Status interface.
func (s Status) Value() interface{} {
	return s.Metadata
}

// Err implements the scanner.Status interface.
func (s Status) Err() error {
	return s.err
}
