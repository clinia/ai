package api

import "net/http"

// SparseEmbeddingOptions represents per-call options for sparse embeddings.
type SparseEmbeddingOptions struct {
	Headers          http.Header
	BaseURL          *string
	ProviderMetadata *ProviderMetadata
}

func (o SparseEmbeddingOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }
