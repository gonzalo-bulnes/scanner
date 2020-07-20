package scanner

import (
	"sort"
	"strings"
	"testing"
)

type Example struct {
	status string
}

func (e Example) Check() Status {
	return e.status
}

func TestScan(t *testing.T) {

	setup := func() (scanner *Scanner, services []Service, output chan Status, done chan bool) {
		scanner = New()
		services = []Service{
			Example{status: "ok"},
			Example{status: "good"},
		}
		output = make(chan Status, len(services))
		done = make(chan bool, 1)
		return
	}

	t.Run("checks all services", func(t *testing.T) {
		scanner, services, output, done := setup()

		scanner.Scan(output, done, services...)
		<-done

		if checks, expected := len(output), len(services); checks != expected {
			t.Errorf("Expected %d elements, got %d", expected, checks)
		}
	})

	t.Run("allows to retrieve the status of all services", func(t *testing.T) {
		scanner, services, output, done := setup()

		scanner.Scan(output, done, services...)
		<-done

		all := []string{}
		for status := range output {
			all = append(all, status.(string))
		}
		sort.Strings(all)

		if result, expected := strings.Join(all, " "), "good ok"; result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}
