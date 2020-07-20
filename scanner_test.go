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

	setup := func() (scanner *Scanner, services []Service) {
		services = []Service{
			Example{status: "ok"},
			Example{status: "good"},
		}
		scanner = New()
		return
	}

	t.Run("checks all services", func(t *testing.T) {
		scanner, services := setup()

		scanner.Scan(services...)

		if size, expected := len(scanner.Output), len(services); size != expected {
			t.Errorf("Expected %d elements, got %d", expected, size)
		}
	})

	t.Run("allows to retrieve the status of all services", func(t *testing.T) {
		scanner, services := setup()

		scanner.Scan(services...)

		all := []string{}
		for status := range scanner.Output {
			all = append(all, status.(string))
		}
		sort.Strings(all)

		if result, expected := strings.Join(all, " "), "good ok"; result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}
