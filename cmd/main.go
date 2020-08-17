package main

import (
	"time"

	"github.com/gonzalo-bulnes/scanner/cmd/cli"
)

func main() {
	cli := cli.New()

	timeout := 20 * time.Second
	cli.Run(timeout)
}
