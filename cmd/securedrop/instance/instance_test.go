package instance

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestInstance_GetMetadata(t *testing.T) {
	t.Run("returns request errors", func(t *testing.T) {
		i := New(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("request error")
		}
		expected := "error querying metadata: request error"

		_, err := i.GetMetadata(context.Background())
		if err == nil {
			t.Fatal("Expected error, got none.")
		}

		if actual := err.Error(); actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("returns error when response status is not HTTP 200 OK", func(t *testing.T) {
		i := New(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
			}, nil
		}
		expected := "instance responded with HTTP 503"

		_, err := i.GetMetadata(context.Background())
		if err == nil {
			t.Fatal("Expected error, got none.")
		}

		if actual := err.Error(); actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("returns error when response body is not JSON", func(t *testing.T) {
		i := New(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			invalidJSON := ""
			body := ioutil.NopCloser(bytes.NewBufferString(invalidJSON))
			return &http.Response{
				Body:       body,
				StatusCode: http.StatusOK,
			}, nil
		}

		expected := "error deserializing instance response: unexpected end of JSON input"

		_, err := i.GetMetadata(context.Background())
		if err == nil {
			t.Fatal("Expected error, got none.")
		}
		if result := err.Error(); result != expected {
			t.Errorf("Expected error to be '%v', was '%v'", expected, result)
		}
	})

	t.Run("known values", func(t *testing.T) {
		response, err := os.Open("testdata/metadata.json")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		i := New(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       ioutil.NopCloser(bufio.NewReader(response)),
				StatusCode: http.StatusOK,
			}, nil
		}

		expectedV2SourceURL := "nyttips4bmquxfzw.onion"

		expected := Metadata{
			AllowDocumentUploads: true,
			GPGFingerprint:       "C0A7BC8D9694BF2FC5EF31EB30614C130E1E864A",
			SecureDropVersion:    "1.5.0",
			ServerOS:             "16.04",
			SupportedLanguages:   []string{"en_US"},
			V2SourceURL:          &expectedV2SourceURL,
		}

		metadata, err := i.GetMetadata(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		checkMetadata(t, expected, metadata)
	})
}

func TestInstance_Check(t *testing.T) {

	t.Run("returns request errors", func(t *testing.T) {
		i := New(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("request error")
		}

		expected := "status check error: error querying metadata: request error"

		s := i.Check(context.Background())
		if s.Err() == nil {
			t.Fatal("Expected error, got none.")
		}

		if actual := s.Err().Error(); actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}

		if s.Value().(Metadata).Available {
			t.Errorf("Expected instance not to be reported as available")
		}
	})

	t.Run("reports the instance availability status and metadata", func(t *testing.T) {
		response, err := os.Open("testdata/metadata.json")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		i := New(nil, "some.onion")
		i.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       ioutil.NopCloser(bufio.NewReader(response)),
				StatusCode: http.StatusOK,
			}, nil
		}

		expectedV2SourceURL := "nyttips4bmquxfzw.onion"

		expected := Metadata{
			Available:            true,
			AllowDocumentUploads: true,
			GPGFingerprint:       "C0A7BC8D9694BF2FC5EF31EB30614C130E1E864A",
			SecureDropVersion:    "1.5.0",
			ServerOS:             "16.04",
			SupportedLanguages:   []string{"en_US"},
			V2SourceURL:          &expectedV2SourceURL,
		}

		s := i.Check(context.Background())
		if err := s.Err(); err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		metadata := s.Value().(Metadata)

		if !metadata.Available {
			t.Errorf("Expected instance to be reported as available")
		}

		checkMetadata(t, expected, metadata)
	})
}

func checkMetadata(t *testing.T, expected, actual Metadata) {
	if a, e := actual.Available, expected.Available; a != e {
		t.Errorf("Expected Available to be %t, got %t", e, a)
	}

	if a, e := actual.AllowDocumentUploads, expected.AllowDocumentUploads; a != e {
		t.Errorf("Expected AllowDocumentUploads to be:\n'%+v', got\n'%+v'\n", e, a)
	}

	if a, e := actual.GPGFingerprint, expected.GPGFingerprint; a != e {
		t.Errorf("Expected GPGFingerprint to be:\n'%+v', got\n'%+v'\n", e, a)
	}

	if a, e := actual.SecureDropVersion, expected.SecureDropVersion; a != e {
		t.Errorf("Expected SecureDropVersion to be:\n'%+v', got\n'%+v'\n", e, a)
	}

	if a, e := actual.ServerOS, expected.ServerOS; a != e {
		t.Errorf("Expected ServerOS to be:\n'%+v', got\n'%+v'\n", e, a)
	}

	if a, e := actual.SupportedLanguages[0], expected.SupportedLanguages[0]; a != e {
		t.Errorf("Expected SupportedLanguages[0] to be:\n'%+v', got\n'%+v'\n", e, a)
	}

	if a, e := *actual.V2SourceURL, *expected.V2SourceURL; a != e {
		t.Errorf("Expected V2SourceURL to be:\n'%+v', got\n'%+v'\n", e, a)
	}

	if url := actual.V3SourceURL; url != nil {
		t.Errorf("Expected V3SourceURL to be nil, got\n'%+v'\n", url)
	}
}
