package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// EmbeddingOption mutates per-call embedding configuration.
type EmbeddingOption func(*api.EmbeddingOptions)

// WithEmbeddingHeaders sets extra HTTP headers for this embedding call.
// Only applies to HTTP-backed providers.
func WithEmbeddingHeaders(headers http.Header) EmbeddingOption {
	return func(o *api.EmbeddingOptions) {
		o.Headers = headers
	}
}

// WithEmbeddingProviderMetadata sets provider-specific metadata for the embedding call.
func WithEmbeddingProviderMetadata(provider string, metadata any) EmbeddingOption {
	return func(o *api.EmbeddingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

// WithEmbeddingBaseURL sets the base URL for the embedding API endpoint.
func WithEmbeddingBaseURL(baseURL string) EmbeddingOption {
	url := baseURL
	return func(o *api.EmbeddingOptions) {
		o.BaseURL = &url
	}
}

// WithEmbeddingEmbeddingOptions sets the entire api.EmbeddingOptions struct.
func WithEmbeddingEmbeddingOptions(embeddingOptions api.EmbeddingOptions) EmbeddingOption {
	return func(o *api.EmbeddingOptions) {
		*o = embeddingOptions
	}
}

// buildEmbeddingConfig combines multiple options into a single api.EmbeddingOptions struct.
func buildEmbeddingConfig(opts []EmbeddingOption) api.EmbeddingOptions {
	config := api.EmbeddingOptions{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
