package securedrop

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gonzalo-bulnes/scanner"
	"github.com/gonzalo-bulnes/scanner/cmd/tor"
)

const instanceMetadataURLPattern = "http://%s/metadata"

// Metadata stores JSON metadata about a SecureDrop instance.
type Metadata struct {
	Fingerprint string `json:"gpg_fpr"`
	Version     string `json:"sd_version"`
}

// Instance represents a SecureDrop instance.
type Instance struct {
	client    *tor.Client
	doRequest func(*http.Request) (*http.Response, error)
	Available bool
	Info      Metadata
	URL       string `json:"Url"`
}

// CSV implements the sdstatus.Information interface.
func (i *Instance) CSV() string {
	return fmt.Sprintf("%s,%s,%s", i.URL, i.Info.Version, i.Info.Fingerprint)
}

// NewInstance returns a new SecureDrop instance.
func NewInstance(client *tor.Client, url string) *Instance {
	i := &Instance{
		client: client,
		URL:    url,
	}
	i.doRequest = client.Do
	return i
}

// Status of a SecureDrop instance.
type Status struct {
	value *Instance
	err   error
}

// Value implements the scanner.Status interface.
func (s Status) Value() interface{} {
	return s.value
}

// Err implements the scanner.Status interface.
func (s Status) Err() error {
	return s.err
}

// Check implements the scanner.Service interface.
func (i *Instance) Check(ctx context.Context) scanner.Status {

	metadataURL := fmt.Sprintf(instanceMetadataURLPattern, i.URL)
	req, err := http.NewRequestWithContext(ctx, "GET", metadataURL, nil)
	if err != nil {
		err = fmt.Errorf("status request creation failed: %w", err)
	}

	resp, err := i.doRequest(req)
	if err != nil {
		return Status{err: fmt.Errorf("status check error: %w", err)}
	}

	if resp.StatusCode != http.StatusOK {
		return Status{err: fmt.Errorf("instance responded with HTTP %d", resp.StatusCode)}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Status{err: fmt.Errorf("error reading status check response: %w", err)}
	}

	err = json.Unmarshal(body, &i.Info)
	if err != nil {
		return Status{err: fmt.Errorf("error deserializing status check response: %w", err)}
	}

	return Status{value: i}
}
