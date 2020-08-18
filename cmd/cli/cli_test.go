package cli

import (
	"context"
	"testing"
	"time"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/securedrop/directory"
)

var ms time.Duration = 1_000_000 // ns

func TestRun(t *testing.T) {

	setup := func() *CLI {
		cli := New()

		// Returns a zero-value configuration.
		cli.configure = func() config {
			return config{}
		}

		// Does nothing.
		cli.getDirectory = func(context.Context) ([]directory.Entry, error) {
			return nil, nil
		}

		// Takes a second, can be cancelled, exits gracefully (closing its output channel).
		cli.scan = func(ctx context.Context, output chan<- scanner.Status, services ...scanner.Service) {
			select {
			case <-ctx.Done():
				close(output)
				return
			case <-time.After(time.Second):
				close(output)
				return
			}
		}
		return cli
	}

	t.Run("times out when the timeout option is set", func(t *testing.T) {
		cli := setup()
		cli.configure = func() config {
			timeout, err := time.ParseDuration("50ms")
			if err != nil {
				t.Fatalf("error setting up test case: %v", err)
			}

			return config{
				timeout: timeout,
			}
		}

		start := time.Now()
		cli.Run()
		end := time.Now()
		elapsed := end.Sub(start)

		if expected := 60 * time.Millisecond; elapsed > expected {
			t.Errorf("Expected checks to take at most %dms (conservative approximation), took %dms", expected/ms, elapsed/ms)
		}
	})
}
