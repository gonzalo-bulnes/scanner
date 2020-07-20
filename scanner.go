package scanner

import "sync"

// Service represents a service which status can be checked.
type Service interface {
	Check() Status
}

// Status represents the status of a service.
type Status interface{}

// CheckFunc is a function that the scanner can use to check the status of a service.
type CheckFunc func() Status

// Scanner allows to scan services checking their status.
type Scanner struct {
	Output chan Status
}

// Scan checks the status of multiple services concurrently.
func (s *Scanner) Scan(services ...Service) {
	s.Output = make(chan Status, len(services))

	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)
		go func(wg *sync.WaitGroup, check CheckFunc) {
			defer wg.Done()
			s.Output <- check()
		}(&wg, service.Check)
	}
	wg.Wait()

	close(s.Output)
}

// New returns a new scanner.
func New() *Scanner {
	return &Scanner{}
}
