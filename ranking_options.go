package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// RankingOption mutates RankingOptions (builder pattern).
type RankingOption func(*api.RankingOptions)

// WithRankingHeaders sets extra HTTP headers for this ranking call.
// Only applies to HTTP-backed providers.
func WithRankingHeaders(headers http.Header) RankingOption {
	return func(o *api.RankingOptions) {
		o.Headers = headers
	}
}

// WithRankingBaseURL sets the base URL for the ranking API endpoint.
func WithRankingBaseURL(baseURL string) RankingOption {
	url := baseURL
	return func(o *api.RankingOptions) {
		o.BaseURL = &url
	}
}

// WithRankingRequestID sets the request identifier associated with the ranking call.
// WithRankingProviderMetadata sets provider-specific metadata for the ranking call.
func WithRankingProviderMetadata(provider string, metadata any) RankingOption {
	return func(o *api.RankingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

func buildRankingConfig(opts []RankingOption) api.RankingOptions {
	config := api.RankingOptions{}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
