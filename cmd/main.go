package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop/directory"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop/instance"
	"github.com/gonzalo-bulnes/scanner/cmd/tor"
)

func main() {
	client, err := tor.NewClient()
	checkError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	d := directory.New(client)
	entries, err := d.Get(ctx)
	checkError(err)

	services := make([]scanner.Service, len(entries))
	for i, entry := range entries {
		services[i] = instance.New(client, entry.OnionAddress)
	}

	s := scanner.New()

	output := make(chan scanner.Status, len(services))
	go s.Scan(ctx, output, services...)

	for status := range output {
		checkError(status.Err())

		metadata := status.Value().(instance.Metadata)
		line := instance.NewOutputLineFromMetadata(metadata)
		bytes, err := line.JSONL()
		checkError(err)

		fmt.Println(string(bytes))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
