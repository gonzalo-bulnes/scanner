package securedrop

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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

func TestInstanceCheck(t *testing.T) {
	t.Run("returns request errors", func(t *testing.T) {
		i := NewInstance(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("request error")
		}
		expected := "status check error: request error"

		result := i.Check(context.Background()).Err().Error()

		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("returns error when response status is not HTTP 200 OK", func(t *testing.T) {
		i := NewInstance(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusForbidden,
			}, nil
		}

		expected := "instance responded with HTTP 403"

		result := i.Check(context.Background()).Err().Error()

		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("returns error when response body is not JSON", func(t *testing.T) {
		i := NewInstance(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			invalidJSON := ""
			body := ioutil.NopCloser(bytes.NewBufferString(invalidJSON))
			return &http.Response{
				Body:       body,
				StatusCode: http.StatusOK,
			}, nil
		}

		expected := "error deserializing status check response: unexpected end of JSON input"

		err := i.Check(context.Background()).Err()
		if err == nil {
			t.Fatal("Expected error, got none.")
		}
		if result := err.Error(); result != expected {
			t.Errorf("Expected error to be '%v', was '%v'", expected, result)
		}
	})

	t.Run("can be serialized after successful check", func(t *testing.T) {
		i := NewInstance(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			body := ioutil.NopCloser(bytes.NewBufferString(`{"hello": "world"}`))
			return &http.Response{
				Body:       body,
				StatusCode: http.StatusOK,
			}, nil
		}

		expectedCSV := "some.onion,,"

		result := i.Check(context.Background())
		if err := result.Err(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if csv := i.CSV(); csv != expectedCSV {
			t.Errorf("Expected '%s', got '%s'", expectedCSV, csv)
		}
	})
}
