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

// WithEmbeddingBaseURL sets the base URL for the embedding API.
func WithEmbeddingBaseURL(baseURL string) EmbeddingOption {
	return func(o *EmbeddingOptions) {
		o.BaseURL = &baseURL
	}
}

// EmbeddingOptions represents the options for generating embeddings.
type EmbeddingOptions struct {
	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header

	// BaseURL is the base URL for the embedding API.
	BaseURL *string
}
