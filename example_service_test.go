package scanner

import (
	"context"
	"fmt"
	"time"
)

// Example of a service which status can be checked.
type Example struct {
	// lets suppose a check takes some time
	duration time.Duration
	// eventually the check returns this status (unless it is cancelled)
	status ExampleStatus
}

// Check is an example check function that supports cancellation.
func (e Example) Check(ctx context.Context) Status {
	select {
	case <-ctx.Done():
		return ExampleStatus{err: fmt.Errorf("check cancelled: %w", ctx.Err())}
	case <-time.After(e.duration):
		return e.status
	}
}

// ExampleStatus is a simple example with a string value.
//
// The value could have any type, and scanner doesn't modify it in any way.
type ExampleStatus struct {
	value string
	err   error
}

// Value implements the Status interface.
func (s ExampleStatus) Value() interface{} {
	return s.value
}

// Err implements the Status interface.
func (s ExampleStatus) Err() error {
	return s.err
}

// Ensures the example is included in the documentation.
func ExampleService() {}
