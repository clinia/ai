package api

import (
	"context"
)

// FileProcessor represents a "one file in -> N files out" processor.
// Returns either a single file payload or an archive (e.g., zip).
type FileProcessor interface {
	// ProviderName returns the name of the provider for logging purposes.
	ProviderName() string

	// ModelID returns the provider-specific model/service ID for logging purposes.
	ModelID() string

	// SupportedUrls allows the SDK to bypass pre-downloading for some URLs (like LanguageModel).
	// If empty, the SDK downloads inputs itself.
	SupportedUrls() []SupportedURL

	// MaxBytesPerCall optionally announces input size limits (nil if unknown).
	MaxBytesPerCall() *int64

	// SupportsParallelCalls indicates whether multiple calls can run in parallel safely.
	SupportsParallelCalls() bool

	// Process uploads a single file and returns the full result in-memory.
	// The result may be a single file or an archive (opaque bytes + media type).
	DoProcess(ctx context.Context, req FileRequest, opts ...FileCallOption) (FileResponse, error)
}

// FileRequest is a generic "one input file" request.
// The SDK will form multipart/form-data and map Params to form fields.
type FileRequest struct {
	// File bytes to upload.
	InlineFile []byte

	// Filename used in multipart (e.g., "document.pdf").
	Filename string

	// MediaType of the input (e.g., "application/pdf").
	MediaType string

	// Params are additional form fields, e.g., {"dpi":"144","pages":"1-3,5"}.
	Params map[string]string
}

// FileResponse represents the full, non-streaming result.
type FileResponse struct {
	// Payload is the complete output bytes (single file OR archive).
	Payload []byte `json:"payload"`

	// MediaType is the MIME type of the payload (e.g., "application/zip", "image/png").
	MediaType string `json:"media_type"`

	// SuggestedFilename from Content-Disposition, if available.
	SuggestedFilename string `json:"suggested_filename,omitzero"`

	// ProviderMetadata can carry provider-specific details (optional).
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`

	// RequestInfo is optional request information for telemetry and debugging purposes.
	RequestInfo *RequestInfo `json:"request,omitzero"`

	// ResponseInfo is optional response information for telemetry and debugging purposes.
	ResponseInfo *ResponseInfo `json:"response,omitzero"`

	// Warnings is a list of warnings that occurred during the call,
	// e.g., unsupported or adjusted params.
	Warnings []CallWarning `json:"warnings,omitempty"`
}

func (r FileResponse) GetProviderMetadata() *ProviderMetadata { return r.ProviderMetadata }
