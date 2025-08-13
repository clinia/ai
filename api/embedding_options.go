package api

import "net/http"

// EmbeddingOption represent the options for generating embeddings.
type EmbeddingOption func(*EmbeddingOptions)

// WithEmbeddingHeaders sets HTTP headers to be sent with the request.
// Only applicable for HTTP-based providers.
func WithEmbeddingHeaders(headers http.Header) EmbeddingOption {
	return func(o *EmbeddingOptions) {
		o.Headers = headers
	}
}

// WithEmbeddingProviderMetadata sets provider-specific metadata.
func WithEmbeddingProviderMetadata(provider string, metadata any) EmbeddingOption {
	return func(o *EmbeddingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = NewProviderMetadata(nil)
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

// EmbeddingOptions represents the options for generating embeddings.
type EmbeddingOptions struct {
	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header

	// ProviderMetadata contains additional provider-specific metadata.
	// The metadata is passed through to the provider from the AI SDK and enables
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata `json:"provider_metadata,omitzero"`
}
