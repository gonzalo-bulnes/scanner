Scanner
=======

Scan the status of multiple services concurrently.

Usage
-----

```go
// main.go

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gonzalo-bulnes/scanner"
)

// example service which status check takes time.
type example struct {
	duration time.Duration
	status   string
	name     string
}

// Check is an example check function that supports cancellation.
func (e example) Check(ctx context.Context) scanner.Status {
	select {
	case <-ctx.Done():
		return exampleStatus{err: fmt.Errorf("%s: %w", e.name, ctx.Err())}
	case <-time.After(e.duration):
		return exampleStatus{value: e.status}
	}
}

// exampleStatus also conveys errors.
type exampleStatus struct {
	value string
	err   error
}

func (s exampleStatus) Value() interface{} {
	return s.value
}

func (s exampleStatus) Err() error {
	return s.err
}

func main() {
	services := []scanner.Service{
		example{name: "fast           ", status: "ok"},
		example{name: "not fast at all", duration: 3 * time.Second, status: "too slow"},
		example{name: "not so fast    ", duration: 500 * time.Millisecond, status: "running slow"},
	}

	s := scanner.New()

	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	output := make(chan scanner.Status, len(services))
	go s.Scan(ctx, output, services...)

	for status := range output {
		if err := status.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(status.Value().(string))
		}
	}
}
```
