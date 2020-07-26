package scanner

import (
	"context"
	"fmt"
	"time"
)

func ExampleScanner_Scan() {

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

	output := make(chan Status, len(services))
	go s.Scan(ctx, output, services...)

	for status := range output {
		if err := status.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(status.Value().(string))
		}
	}

	// Output (in any order because of concurrency):
	// unavailable service: responded with HTTP 503 Unavailable
	// fast service: ok
	// not so fast service: running slow
	// check cancelled: context deadline exceeded
}
