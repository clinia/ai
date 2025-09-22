package api

import "net/http"

// ChunkingOptions represents per-call options for chunking models.
type ChunkingOptions struct {
	Headers          http.Header
	BaseURL          *string
	ProviderMetadata *ProviderMetadata
}

func (o ChunkingOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }
