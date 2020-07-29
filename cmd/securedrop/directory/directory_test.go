package directory

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

func TestDirectory_Get(t *testing.T) {
	t.Run("returns request errors", func(t *testing.T) {
		d := New(nil)
		d.doRequest = func(*http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("request error")
		}
		expected := "error querying directory: request error"

		_, err := d.Get(context.Background())
		if err == nil {
			t.Fatal("Expected error, got none.")
		}

		if actual := err.Error(); actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("returns error when response status is not HTTP 200 OK", func(t *testing.T) {
		d := New(nil)
		d.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
			}, nil
		}
		expected := "directory responded with HTTP 503"

		_, err := d.Get(context.Background())
		if err == nil {
			t.Fatal("Expected error, got none.")
		}

		if actual := err.Error(); actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("returns error when response body is not JSON", func(t *testing.T) {
		d := New(nil)
		d.doRequest = func(*http.Request) (*http.Response, error) {
			invalidJSON := ""
			body := ioutil.NopCloser(bytes.NewBufferString(invalidJSON))
			return &http.Response{
				Body:       body,
				StatusCode: http.StatusOK,
			}, nil
		}

		expected := "error deserializing directory response: unexpected end of JSON input"

		_, err := d.Get(context.Background())
		if err == nil {
			t.Fatal("Expected error, got none.")
		}
		if result := err.Error(); result != expected {
			t.Errorf("Expected error to be '%v', was '%v'", expected, result)
		}
	})

	t.Run("known values", func(t *testing.T) {
		response, err := os.Open("testdata/response.json")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		d := New(nil)
		d.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       ioutil.NopCloser(bufio.NewReader(response)),
				StatusCode: http.StatusOK,
			}, nil
		}

		expected := []Entry{
			{
				Title:            "Public Intelligence",
				Slug:             "public-intelligence",
				DirectoryURL:     "https://securedrop.org/directory/public-intelligence/",
				FirstPublishedAt: "2018-03-28T05:39:21.461705Z",
				LandingPageURL:   "https://publicintelligence.net/contribute/",
				OnionAddress:     "arujlhu2zjjhc3bw.onion",
				OrganisationLogo: OrganisationLogo{
					URL:    "/media/images/pilogo.max-1500x1500.png",
					Width:  500,
					Height: 450,
				},
				OrganisationDescription: "Public Intelligence is a collaborative project to defend the public's right to information.",
				Languages:               []string{"English"},
				Topics:                  []string{"business", "environment", "government", "health", "war"},
				Countries:               []string{"All countries"},
			},
			{
				Title:            "The Guardian",
				Slug:             "guardian",
				DirectoryURL:     "https://securedrop.org/directory/guardian/",
				FirstPublishedAt: "2017-11-09T19:57:47.051165Z",
				LandingPageURL:   "https://www.theguardian.com/securedrop",
				OnionAddress:     "33y6fjyhs3phzfjj.onion",
				OrganisationLogo: OrganisationLogo{
					URL:    "/media/images/Guardian_titlepiece.max-1500x1500.png",
					Width:  1500,
					Height: 481,
				},
				OrganisationDescription: "The Guardian is a British daily newspaper.",
				Languages:               []string{"All languages", "English"},
				Topics: []string{"business", "civil liberties", "environment",
					"government", "health", "national security", "technology"},
				Countries: []string{"All countries"},
			},
		}

		entries, err := d.Get(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if n, m := len(entries), len(expected); n != m {
			t.Fatalf("Expected %d  entries, got %d", m, n)
		}

		for i, entry := range entries {
			checkEntry(t, expected[i], entry, i)
		}
	})
}

func checkEntry(t *testing.T, expected, entry Entry, id int) {
	if a, e := entry.Title, expected.Title; a != e {
		t.Errorf("Expected entry[%d] Title to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.Slug, expected.Slug; a != e {
		t.Errorf("Expected entry[%d] Slug to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.DirectoryURL, expected.DirectoryURL; a != e {
		t.Errorf("Expected entry[%d] DirectoryURL to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.FirstPublishedAt, expected.FirstPublishedAt; a != e {
		t.Errorf("Expected entry[%d] FirstPublishedAt to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.LandingPageURL, expected.LandingPageURL; a != e {
		t.Errorf("Expected entry[%d] LandingPageURL to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.OnionAddress, expected.OnionAddress; a != e {
		t.Errorf("Expected entry[%d] OnionAddress to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.OrganisationLogo, expected.OrganisationLogo; a != e {
		t.Errorf("Expected entry[%d] OrganisationLogo to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.OrganisationDescription, expected.OrganisationDescription; a != e {
		t.Errorf("Expected entry[%d] OrganisationDescription to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.Languages[0], expected.Languages[0]; a != e {
		t.Errorf("Expected entry[%d] Languages to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.Topics[1], expected.Topics[1]; a != e {
		t.Errorf("Expected entry[%d] Topics to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
	if a, e := entry.Countries[0], expected.Countries[0]; a != e {
		t.Errorf("Expected entry[%d] Countries to be:\n'%+v', got\n'%+v'\n", id, e, a)
	}
}
