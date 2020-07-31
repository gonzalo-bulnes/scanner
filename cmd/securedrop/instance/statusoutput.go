package instance

import (
	"encoding/json"
	"fmt"
)

// StatusOutput allows to output relevant status information about a collection of SecureDrop instances.
type StatusOutput []StatusOutputLine

// NewOutputFromMetadata returns relevant metadata for output purposes.
func NewOutputFromMetadata(m []Metadata) StatusOutput {
	output := make(StatusOutput, len(m))
	for i := 0; i < len(m); i++ {
		output[i] = NewOutputLineFromMetadata(m[i])
	}

	return output
}

// JSON produces a JSON output.
func (o StatusOutput) JSON() ([]byte, error) {
	return json.Marshal(o)
}

// CSV produces a CSV output.
func (o StatusOutput) CSV() (csv string) {
	for _, line := range o {
		csv += line.CSV()
	}
	return
}

// StatusOutputLine represents the status information of a single SecureDrop instance for output purposes.
type StatusOutputLine struct {
	Info      StatusOutputInfo `json:"Info"`
	URL       string           `json:"Url"`
	Available bool             `json:"Available"`
}

// NewOutputLineFromMetadata returns relevant metadata for output purposes.
func NewOutputLineFromMetadata(m Metadata) (line StatusOutputLine) {
	line.Info = StatusOutputInfo{
		SecureDropVersion: m.SecureDropVersion,
		GPGFingerprint:    m.GPGFingerprint,
	}
	line.Available = m.Available
	line.URL = m.URL

	return
}

// CSV produces a CSV output.
func (l StatusOutputLine) CSV() (csv string) {
	csv += fmt.Sprintf("%s,%s,%s\n", l.URL, l.Info.SecureDropVersion, l.Info.GPGFingerprint)

	return
}

// JSONL produces a JSON Lines output.
func (l StatusOutputLine) JSONL() ([]byte, error) {
	return json.Marshal(l)
}

// StatusOutputInfo is only defined for JSON marshalling purposes.
type StatusOutputInfo struct {
	SecureDropVersion string `json:"sd_version"`
	GPGFingerprint    string `json:"gpg_fpr"`
}
