package scanner

import (
	"context"
	"fmt"
	"sort"
	"time"
)

func ExampleScanner_ScanAndWait() {

	// Example and ExampleStatus are defined in the Service type example.

	services := []Service{
		Example{
			status: ExampleStatus{
				value: "fast service: ok",
			},
		},
		Example{
			duration: 200 * time.Millisecond,
			status: ExampleStatus{
				value: "not so fast service: running slow",
			},
		},
		Example{
			duration: 800 * time.Second,
			status: ExampleStatus{
				value: "slow service: too slow",
			},
		},
		Example{
			duration: 400 * time.Millisecond,
			status: ExampleStatus{
				err: fmt.Errorf("unavailable service: responded with HTTP 503 Unavailable"),
			},
		},
	}

	s := New()

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	responses := s.ScanAndWait(ctx, services...)

	// For testing purposes, do some sorting to ensure a stable output:
	var line string
	lines := []string{}
	for _, status := range responses {
		if err := status.Err(); err != nil {
			line = err.Error()
		} else {
			line = status.Value().(string)
		}
		lines = append(lines, line)
	}
	sort.Strings(lines)

	for _, line := range lines {
		fmt.Println(line)
	}

	// Output:
	// check cancelled: context deadline exceeded
	// fast service: ok
	// not so fast service: running slow
	// unavailable service: responded with HTTP 503 Unavailable
}
