package instance

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestStatusOuput_JSON(t *testing.T) {
	t.Run("known values", func(t *testing.T) {

		output := StatusOutput{
			{
				Info: StatusOutputInfo{
					SecureDropVersion: "0.6",
					GPGFingerprint:    "3392A1CE68FE779A95FCAF04EDA0FB6F53FA9093",
				},
				URL:       "m4hynbhhctdk27jr.onion",
				Available: true,
			},
			{
				Info: StatusOutputInfo{
					SecureDropVersion: "0.6",
					GPGFingerprint:    "7C24A77EED0D50838E3315BD7A38590B2996F0C2",
				},
				URL:       "ftugftwajmgsmoau.onion",
				Available: true,
			},
		}

		f, err := os.Open("testdata/output.json")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		expected, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		actual, err := output.JSON()
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if normalizeJSON(actual) != normalizeJSON(expected) {
			t.Errorf("Expected:\n%s, got\n%s\n", normalizeJSON(expected), normalizeJSON(actual))
		}
	})
}

func normalizeJSON(sample []byte) (s string) {
	s = string(sample)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return
}

func TestStatusOuput_CSV(t *testing.T) {
	t.Run("known values", func(t *testing.T) {

		output := StatusOutput{
			{
				Info: StatusOutputInfo{
					SecureDropVersion: "0.6",
					GPGFingerprint:    "3392A1CE68FE779A95FCAF04EDA0FB6F53FA9093",
				},
				URL:       "m4hynbhhctdk27jr.onion",
				Available: true,
			},
			{
				Info: StatusOutputInfo{
					SecureDropVersion: "0.6",
					GPGFingerprint:    "7C24A77EED0D50838E3315BD7A38590B2996F0C2",
				},
				URL:       "ftugftwajmgsmoau.onion",
				Available: true,
			},
		}

		f, err := os.Open("testdata/output.csv")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		expected, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		actual := output.CSV()

		if actual != string(expected) {
			t.Errorf("Expected:\n%s, got\n%s\n", string(expected), actual)
		}
	})
}
