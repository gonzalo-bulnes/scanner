package securedrop

import (
	"testing"
)

func TestInstanceCSV(t *testing.T) {
	t.Run("known values", func(t *testing.T) {
		i := Instance{
			Info: Metadata{
				Fingerprint: "F6E0E2901B787C3721E1C0BF4BD6284B525A3DF4",
				Version:     "1.4.1",
			},
			URL: "zdf4nikyuswdzbt6.onion",
		}
		expected := "zdf4nikyuswdzbt6.onion,1.4.1,F6E0E2901B787C3721E1C0BF4BD6284B525A3DF4"

		if result := i.CSV(); result != expected {
			t.Errorf("Expected '%s', got %s", expected, result)
		}
	})
}
