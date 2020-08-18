package scanner

import (
	"context"
	"sort"
	"strings"
	"testing"
	"time"
)

var ms time.Duration = 1_000_000 // ns

func setup() (scanner *Scanner, services []Service, output chan Status) {
	scanner = New()
	services = []Service{
		Example{status: ExampleStatus{value: "ok"}},
		Example{status: ExampleStatus{value: "good"}},
	}
	output = make(chan Status, len(services))
	return
}

func TestScan(t *testing.T) {

	t.Run("checks all services", func(t *testing.T) {
		scanner, services, output := setup()

		scanner.Scan(context.Background(), output, services...)

		responses := []Status{}
		for response := range output {
			responses = append(responses, response)
		}

		if checks, expected := len(responses), len(services); checks != expected {
			t.Errorf("Expected %d elements, got %d", expected, checks)
		}
	})

	t.Run("allows to retrieve the status of all services", func(t *testing.T) {
		scanner, services, output := setup()

		scanner.Scan(context.Background(), output, services...)

		all := []string{}
		for status := range output {
			all = append(all, status.Value().(string))
		}
		sort.Strings(all)

		if result, expected := strings.Join(all, " "), "good ok"; result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

func TestScanAndWait(t *testing.T) {

	t.Run("checks all services", func(t *testing.T) {
		scanner, services, _ := setup()

		responses := scanner.ScanAndWait(context.Background(), services...)

		if checks, expected := len(responses), len(services); checks != expected {
			t.Errorf("Expected %d elements, got %d", expected, checks)
		}
	})

	t.Run("allows to retrieve the status of all services", func(t *testing.T) {
		scanner, services, _ := setup()

		responses := scanner.ScanAndWait(context.Background(), services...)

		all := []string{}
		for _, status := range responses {
			all = append(all, status.Value().(string))
		}
		sort.Strings(all)

		if result, expected := strings.Join(all, " "), "good ok"; result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

func TestSetWorkerCount(t *testing.T) {

	t.Run("allows to set a limit to concurrency", func(t *testing.T) {
		scanner, _, _ := setup()
		scanner.SetWorkerCount(1)

		services := []Service{
			Example{status: ExampleStatus{value: "ok"}, duration: 200 * time.Millisecond},
			Example{status: ExampleStatus{value: "good"}, duration: 300 * time.Millisecond},
		}

		start := time.Now()
		_ = scanner.ScanAndWait(context.Background(), services...)
		end := time.Now()
		elapsed := end.Sub(start)

		if expected := 500 * time.Millisecond; elapsed < expected {
			t.Errorf("Expected checks to take at least %dms, took %dms", expected/ms, elapsed/ms)
		}
	})

	t.Run("allows to disable the concurrency limit", func(t *testing.T) {
		scanner, _, _ := setup()
		scanner.SetWorkerCount(1)
		scanner.SetWorkerCount(0)

		services := []Service{
			Example{status: ExampleStatus{value: "ok"}, duration: 200 * time.Millisecond},
			Example{status: ExampleStatus{value: "good"}, duration: 300 * time.Millisecond},
		}

		start := time.Now()
		_ = scanner.ScanAndWait(context.Background(), services...)
		end := time.Now()
		elapsed := end.Sub(start)

		if expected := 400 * time.Millisecond; elapsed > expected {
			t.Errorf("Expected checks to take at most %dms (conservative approximation), took %dms", expected/ms, elapsed/ms)
		}
	})
}
