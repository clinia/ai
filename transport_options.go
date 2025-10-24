package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// TransportOption mutates per-call transport configuration.
type TransportOption func(*api.TransportOptions)

// WithTransportHeaders sets extra HTTP headers for this call.
// Only applies to HTTP-backed providers.
func WithTransportHeaders(headers http.Header) TransportOption {
	return func(o *api.TransportOptions) {
		o.Headers = headers
	}
}

// WithTransportAPIKey sets the API key for this call.
// Only applies to HTTP-backed providers.
func WithTransportAPIKey(apiKey string) TransportOption {
	return func(o *api.TransportOptions) {
		o.APIKey = apiKey
	}
}

// WithTransportProviderMetadata sets provider-specific metadata for the call.
func WithTransportProviderMetadata(provider string, metadata any) TransportOption {
	return func(o *api.TransportOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

// WithTransportBaseURL sets the base URL for the API endpoint.
func WithTransportBaseURL(baseURL string) TransportOption {
	url := baseURL
	return func(o *api.TransportOptions) {
		o.BaseURL = &url
	}
}

// WithTransportUseRawBaseURL instructs HTTP-backed providers to use the provided
// BaseURL as the full request URL without appending a path.
func WithTransportUseRawBaseURL() TransportOption {
	return func(o *api.TransportOptions) {
		o.UseRawBaseURL = true
	}
}

// buildTransportConfig combines multiple options into a single api.TransportOptions struct.
func buildTransportConfig(opts []TransportOption) api.TransportOptions {
	config := api.TransportOptions{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
