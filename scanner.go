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

// Scan checks the status of multiple services concurrently.
func Scan(out chan<- Status, done chan<- bool, services ...Service) {
	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)
		go func(wg *sync.WaitGroup, output chan<- Status, check CheckFunc) {
			defer wg.Done()
			output <- check()
		}(&wg, out, service.Check)
	}
	wg.Wait()

	close(out)
	done <- true
}
