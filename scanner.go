// Package scanner provides primitives to check the status of multiple services concurrently.
package scanner

import (
	"context"
	"sync"
)

// Service represents a service which status can be checked.
type Service interface {
	Check(ctx context.Context) Status
}

// Status represents the status of a service.
type Status interface {
	Value() interface{}
	Err() error
}

// CheckFunc is a function that the scanner can use to check the status of a service.
//
// Ideally, a check function would support being cancelled through its context.
type CheckFunc func(ctx context.Context) Status

// Scanner allows to scan services checking their status.
type Scanner struct{}

// New returns a new scanner.
func New() *Scanner {
	return &Scanner{}
}

// Scan checks the status of multiple services concurrently.
//
// The output channel will be closed when the scan is done.
func (s *Scanner) Scan(ctx context.Context, output chan<- Status, services ...Service) {
	scan(ctx, output, services...)
	close(output)
}

// ScanAndWait checks the status of multiple services concurrently, and returns
// their responses once all the checks are done.
//
// If you want to avoid writing concurrent code, you may still benefit from
// the increased speed of concurrent checks by using this method.
func (s *Scanner) ScanAndWait(ctx context.Context, services ...Service) []Status {
	output := make(chan Status, len(services))
	scan(ctx, output, services...)
	close(output)

	responses := []Status{}
	for status := range output {
		responses = append(responses, status)
	}
	return responses
}

func scan(ctx context.Context, output chan<- Status, services ...Service) {
	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)
		go func(check CheckFunc) {
			defer wg.Done()
			output <- check(ctx)
		}(service.Check)
	}
	wg.Wait()
}
