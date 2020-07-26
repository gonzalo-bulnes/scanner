package scanner

import (
	"context"
	"sort"
	"strings"
	"testing"
)

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
