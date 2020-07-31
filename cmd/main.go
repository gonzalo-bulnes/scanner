package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop/instance"
	"github.com/gonzalo-bulnes/scanner/cmd/tor"
)

func main() {
	client, err := tor.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	services := []scanner.Service{
		instance.New(client, "missing.onion"),
	}

	s := scanner.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output := make(chan scanner.Status, len(services))
	go s.Scan(ctx, output, services...)

	for status := range output {
		if err := status.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(status.Value().(instance.Metadata))
		}
	}
}
