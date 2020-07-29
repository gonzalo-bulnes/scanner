// Package directory provides primitives to interact with the SecureDrop directory.
package directory

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// The SecureDrop directory URL.
const url = "https://securedrop.org/api/v1/directory"

// Entry represents an SecureDrop instance in the directory response.
type Entry struct {
	Title                   string `json:"title"`
	Slug                    string `json:"slug"`
	DirectoryURL            string `json:"directory_url"`
	FirstPublishedAt        string `json:"first_published_at"`
	LandingPageURL          string `json:"landing_page_url"`
	OnionAddress            string `json:"onion_address"`
	OrganisationLogo        `json:"organization_logo"`
	OrganisationDescription string   `json:"organization_description"`
	Languages               []string `json:"languages"`
	Topics                  []string `json:"topics"`
	Countries               []string `json:"countries"`
}

// OrganisationLogo represents the logo of an organisation that operates a SecureDrop instance.
type OrganisationLogo struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Directory allows to query the SecureDrop directory.
type Directory struct {
	client    *http.Client
	doRequest func(*http.Request) (*http.Response, error)
}

// New returns a directory client.
func New(client *http.Client) *Directory {
	d := &Directory{
		client: client,
	}
	d.doRequest = client.Do
	return d
}

// Get fetches directory entries.
func (d *Directory) Get(ctx context.Context) (entries []Entry, err error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		err = fmt.Errorf("status request creation failed: %w", err)
	}

	resp, err := d.doRequest(req)
	if err != nil {
		err = fmt.Errorf("error querying directory: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("directory responded with HTTP %d", resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading directory response: %w", err)
		return
	}

	err = json.Unmarshal(body, &entries)
	if err != nil {
		err = fmt.Errorf("error deserializing directory response: %w", err)
		return
	}
	return
}
