package api

import "net/http"

// EmbeddingOption represent the options for generating embeddings.
type EmbeddingOption func(*EmbeddingOptions)

// EmbeddingOptions represents the options for generating embeddings.
type EmbeddingOptions struct {
	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header

	// APIKey is the API key to be used for authentication.
	// Only applicable for HTTP-based providers that use API key authentication.
	APIKey string

	// BaseURL is the base URL for the API endpoint.
	BaseURL *string

	// UseRawBaseURL, when true, instructs HTTP-backed providers to use the
	// provided BaseURL as the full request URL without appending an API path.
	// This is useful for gateways or endpoints where the BaseURL already
	// contains the complete path.
	// Only applicable for HTTP-based providers that support it.
	UseRawBaseURL bool

	// ProviderMetadata contains additional provider-specific metadata.
	// The metadata is passed through to the provider from the AI SDK and enables
	// provider-specific functionality that can be fully encapsulated in the provider.
	ProviderMetadata *ProviderMetadata
}

func (o EmbeddingOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }
