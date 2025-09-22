package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// ChunkingOptions bundles per-call chunking options for the AI helpers.
type ChunkingOptions struct {
	ChunkingOptions api.ChunkingOptions
}

// ChunkingOption mutates ChunkingOptions (builder pattern).
type ChunkingOption func(*ChunkingOptions)

// WithChunkingHeaders sets extra HTTP headers for this chunking call.
// Only applies to HTTP-backed providers.
func WithChunkingHeaders(headers http.Header) ChunkingOption {
	return func(o *ChunkingOptions) {
		o.ChunkingOptions.Headers = headers
	}
}

// WithChunkingBaseURL sets the base URL for the chunking API endpoint.
func WithChunkingBaseURL(baseURL string) ChunkingOption {
	url := baseURL
	return func(o *ChunkingOptions) {
		o.ChunkingOptions.BaseURL = &url
	}
}

// WithChunkingOptions sets the entire api.ChunkingOptions struct.
func WithChunkingOptions(chunkingOptions api.ChunkingOptions) ChunkingOption {
	return func(o *ChunkingOptions) {
		o.ChunkingOptions = chunkingOptions
	}
}

// WithChunkingRequestID sets the request identifier for the chunking call.
// WithChunkingProviderMetadata sets provider-specific metadata for the chunking call.
func WithChunkingProviderMetadata(provider string, metadata any) ChunkingOption {
	return func(o *ChunkingOptions) {
		if o.ChunkingOptions.ProviderMetadata == nil {
			o.ChunkingOptions.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ChunkingOptions.ProviderMetadata.Set(provider, metadata)
	}
}

func buildChunkingConfig(opts []ChunkingOption) ChunkingOptions {
	config := ChunkingOptions{ChunkingOptions: api.ChunkingOptions{}}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
