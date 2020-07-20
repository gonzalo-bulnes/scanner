Scanner
=======

Scan the status of multiple services concurrently.

Usage
-----

```go
// main.go

package main

import (
	"fmt"
	"time"

	"github.com/gonzalo-bulnes/scanner"
)

// example of a Service
type example struct {
	delay  time.Duration
	status string
}

func (e example) Check() scanner.Status {
	time.Sleep(e.delay)
	return e.status
}

func main() {
	services := []scanner.Service{
		example{
			status: "ok"
	},
		example{
			delay: 2 * time.Second,
			status: "good"
		},
	}

	// create a scanner
	s := scanner.New()

	// start a scan
	output := make(chan scanner.Status, len(services))
	done := make(chan bool, 1)
	go s.Scan(output, done, services...)

	// print the output as a stream
	for status := range output {
		fmt.Println(status)
	}
	<-done

	// alternatively wait until the scan is done
	<-done
	for status := range output {
		fmt.Println(status)
	}
}
```
