package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop"
	"github.com/gonzalo-bulnes/scanner/cmd/tor"
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
	client, err := tor.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	services := []scanner.Service{
		securedrop.NewInstance(client, "missing.onion"),
	}

	s := scanner.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output := make(chan scanner.Status, len(services))
	done := make(chan bool, 1)
	go s.Scan(ctx, output, done, services...)

	for status := range output {
		if err := status.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(status.Value().(*securedrop.Instance).CSV())
		}
	}
	<-done
}
