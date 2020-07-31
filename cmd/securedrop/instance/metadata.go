package instance

// Metadata represents the metadata of a SecureDrop instance.
type Metadata struct {
	URL                  string
	Available            bool
	AllowDocumentUploads bool     `json:"allow_document_uploads"`
	GPGFingerprint       string   `json:"gpg_fpr"`
	SecureDropVersion    string   `json:"sd_version"`
	ServerOS             string   `json:"server_os"`
	SupportedLanguages   []string `json:"supported_languages"`
	V2SourceURL          *string  `json:"v2_source_url"`
	V3SourceURL          *string  `json:"v3_source_url"`
}
