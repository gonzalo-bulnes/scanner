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
type Status interface{}

// CheckFunc is a function that the scanner can use to check the status of a service.
type CheckFunc func(ctx context.Context) Status

// Scanner allows to scan services checking their status.
type Scanner struct{}

// Scan checks the status of multiple services concurrently.
func (s *Scanner) Scan(ctx context.Context, out chan<- Status, done chan<- bool, services ...Service) {
	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, output chan<- Status, check CheckFunc) {
			defer wg.Done()
			output <- check(ctx)
		}(ctx, &wg, out, service.Check)
	}
	wg.Wait()

	close(out)
	done <- true
}

// New returns a new scanner.
func New() *Scanner {
	return &Scanner{}
}
