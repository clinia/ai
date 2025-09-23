package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// ChunkingOption mutates ChunkingOptions (builder pattern).
type ChunkingOption func(*api.ChunkingOptions)

// WithChunkingHeaders sets extra HTTP headers for this chunking call.
// Only applies to HTTP-backed providers.
func WithChunkingHeaders(headers http.Header) ChunkingOption {
	return func(o *api.ChunkingOptions) {
		o.Headers = headers
	}
}

// WithChunkingBaseURL sets the base URL for the chunking API endpoint.
func WithChunkingBaseURL(baseURL string) ChunkingOption {
	url := baseURL
	return func(o *api.ChunkingOptions) {
		o.BaseURL = &url
	}
}

// WithChunkingRequestID sets the request identifier for the chunking call.
// WithChunkingProviderMetadata sets provider-specific metadata for the chunking call.
func WithChunkingProviderMetadata(provider string, metadata any) ChunkingOption {
	return func(o *api.ChunkingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

func buildChunkingConfig(opts []ChunkingOption) api.ChunkingOptions {
	config := api.ChunkingOptions{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
