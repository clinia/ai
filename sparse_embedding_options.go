package ai

import (
	"net/http"

	"go.jetify.com/ai/api"
)

type SparseEmbeddingOption func(*api.SparseEmbeddingOptions)

func WithSparseEmbeddingHeaders(headers http.Header) SparseEmbeddingOption {
	return func(o *api.SparseEmbeddingOptions) { o.Headers = headers }
}

func WithSparseEmbeddingBaseURL(baseURL string) SparseEmbeddingOption {
	return func(o *api.SparseEmbeddingOptions) { u := baseURL; o.BaseURL = &u }
}

func WithSparseEmbeddingProviderMetadata(provider string, metadata any) SparseEmbeddingOption {
	return func(o *api.SparseEmbeddingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(map[string]any{})
		}
		o.ProviderMetadata.Set(provider, metadata)
	}
}

func buildSparseEmbeddingConfig(opts []SparseEmbeddingOption) api.SparseEmbeddingOptions {
	cfg := api.SparseEmbeddingOptions{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
