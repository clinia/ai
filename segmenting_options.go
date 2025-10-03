package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// SegmentingOption mutates per-call segmenting configuration.
type SegmentingOption func(*api.SegmentingOptions)

// WithSegmentingHeaders sets extra HTTP headers for this segmenting call.
func WithSegmentingHeaders(headers http.Header) SegmentingOption {
	return func(o *api.SegmentingOptions) { o.Headers = headers }
}

// WithSegmentingBaseURL sets the base URL for this segmenting call.
func WithSegmentingBaseURL(baseURL string) SegmentingOption {
	return func(o *api.SegmentingOptions) { u := baseURL; o.BaseURL = &u }
}

// WithSegmentingProviderMetadata sets provider-specific metadata for the call.
func WithSegmentingProviderMetadata(provider string, metadata any) SegmentingOption {
	return func(o *api.SegmentingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

func buildSegmentingConfig(opts []SegmentingOption) api.SegmentingOptions {
	cfg := api.SegmentingOptions{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
