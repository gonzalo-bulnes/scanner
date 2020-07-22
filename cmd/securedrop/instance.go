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
	Available bool
	Info      Metadata
	URL       string `json:"Url"`
}

// CSV implements the sdstatus.Information interface.
func (i *Instance) CSV() string {
	return fmt.Sprintf("%s,%s,%s", i.URL, i.Info.Version, i.Info.Fingerprint)
}

// NewInstance returns a new SecureDrop instance.
func NewInstance(url string) *Instance {
	return &Instance{
		URL: url,
	}
}

// Status of a SecureDrop instance.
type Status struct {
	value *Instance
	err   error
}

func (s Status) Value() interface{} {
	return s.value
}

func (s Status) Err() error {
	return s.err
}

func (i *Instance) Check(ctx context.Context) scanner.Status {

	metadataURL := fmt.Sprintf(instanceMetadataURLPattern, i.URL)
	req, err := http.NewRequestWithContext(ctx, "GET", metadataURL, nil)
	if err != nil {
		err = fmt.Errorf("status request creation failed: %w", err)
	}

	c, err := tor.NewClient()
	if err != nil {
		return Status{err: fmt.Errorf("status check error: %w", err)}
	}

	resp, err := c.Do(req)
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
