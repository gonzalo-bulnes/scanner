// Package instance provides a representation of a SecureDrop instance, which status can be checked.
package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/tor"
)

const metadataURLPattern = "http://%s/metadata"

// Instance represents a SecureDrop instance.
type Instance struct {
	client    *tor.Client
	doRequest func(*http.Request) (*http.Response, error)
	URL       string
}

// New returns a new SecureDrop instance.
func New(client *tor.Client, url string) *Instance {
	i := &Instance{
		client: client,
		URL:    url,
	}
	i.doRequest = client.Do
	return i
}

// GetMetadata fetches metadata from a SecureDrop instance.
func (i *Instance) GetMetadata(ctx context.Context) (m Metadata, err error) {

	url := fmt.Sprintf(metadataURLPattern, i.URL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		err = fmt.Errorf("metadata request creation failed: %w", err)
	}

	resp, err := i.doRequest(req)
	if err != nil {
		err = fmt.Errorf("error querying metadata: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("instance responded with HTTP %d", resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading instance response: %w", err)
		return
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		err = fmt.Errorf("error deserializing instance response: %w", err)
		return
	}
	m.URL = i.URL
	return
}

// Check implements the scanner.Service interface.
func (i *Instance) Check(ctx context.Context) scanner.Status {

	metadata, err := i.GetMetadata(ctx)
	if err != nil {
		return Status{err: fmt.Errorf("status check error: %w", err)}
	}
	metadata.Available = true

	return Status{
		Metadata: metadata,
	}
}
