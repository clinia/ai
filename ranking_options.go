package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

// RankingOptions bundles per-call ranking options for the AI helpers.
type RankingOptions struct {
	RankingOptions api.RankingOptions
}

// RankingOption mutates RankingOptions (builder pattern).
type RankingOption func(*RankingOptions)

// WithRankingHeaders sets extra HTTP headers for this ranking call.
// Only applies to HTTP-backed providers.
func WithRankingHeaders(headers http.Header) RankingOption {
	return func(o *RankingOptions) {
		o.RankingOptions.Headers = headers
	}
}

// WithRankingBaseURL sets the base URL for the ranking API endpoint.
func WithRankingBaseURL(baseURL string) RankingOption {
	url := baseURL
	return func(o *RankingOptions) {
		o.RankingOptions.BaseURL = &url
	}
}

// WithRankingOptions sets the entire api.RankingOptions struct.
func WithRankingOptions(rankingOptions api.RankingOptions) RankingOption {
	return func(o *RankingOptions) {
		o.RankingOptions = rankingOptions
	}
}

// WithRankingRequestID sets the request identifier associated with the ranking call.
// WithRankingProviderMetadata sets provider-specific metadata for the ranking call.
func WithRankingProviderMetadata(provider string, metadata any) RankingOption {
	return func(o *RankingOptions) {
		if o.RankingOptions.ProviderMetadata == nil {
			o.RankingOptions.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.RankingOptions.ProviderMetadata.Set(provider, metadata)
	}
}

func buildRankingConfig(opts []RankingOption) RankingOptions {
	config := RankingOptions{RankingOptions: api.RankingOptions{}}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
