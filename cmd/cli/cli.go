package cli

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop/directory"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop/instance"
	"github.com/gonzalo-bulnes/scanner/cmd/tor"
)

// CLI provides a command-line interface for checking the availability of SecireDrop instances.
type CLI struct {
	client       *http.Client
	directory    *directory.Directory
	err          *log.Logger
	out          *log.Logger
	scanner      *scanner.Scanner
	configure    func() config
	getDirectory func(context.Context) ([]directory.Entry, error)
	scan         func(context.Context, chan<- scanner.Status, ...scanner.Service)
}

// New returns a new CLI, or exits if the connection to the Tor network fails.
func New() *CLI {

	cli := &CLI{
		err:       log.New(os.Stderr, "", 0),
		out:       log.New(os.Stdout, "", 0),
		scanner:   scanner.New(),
		configure: makeConfigFromFlags,
	}
	cli.scanner.SetWorkerCount(3)
	cli.scan = cli.scanner.Scan

	client, err := tor.NewClient()
	if err != nil {
		cli.err.Fatalf("Error initialising the command-line interface: %v\n", err)
	}
	cli.client = client

	cli.directory = directory.New(client)
	cli.getDirectory = cli.directory.Get

	return cli
}

// Run checks the availability of all the instances listed in the SecureDrop directory.
func (cli *CLI) Run() {
	cfg := cli.configure()

	ctx := context.Background()
	if cfg.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.timeout)
		defer cancel()
	}

	entries, err := cli.getDirectory(ctx)
	if err != nil {
		cli.err.Fatalf("Error fetching SecureDrop instances list: %v\n", err)
	}

	services := make([]scanner.Service, len(entries))
	for i, entry := range entries {
		services[i] = instance.New(cli.client, entry.OnionAddress)
	}

	output := make(chan scanner.Status, len(services))
	go cli.scan(ctx, output, services...)

	for status := range output {
		if err := status.Err(); err != nil {
			cli.err.Printf("Error checking status of SecureDrop instance: %v\n", err)
			continue
		}

		metadata := status.Value().(instance.Metadata)
		line := instance.NewOutputLineFromMetadata(metadata)
		bytes, err := line.JSONL()
		if err != nil {
			cli.err.Printf("Error formatting status of SecureDrop instance: %v\n", err)
			continue
		}

		cli.out.Printf(string(bytes))
	}
}
