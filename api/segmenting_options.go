package api

import "net/http"

// SegmentingOptions represents per-call options for segmenting models.
// Starts minimal; can grow to include tokenizer, head/tail tokens, etc.
type SegmentingOptions struct {
	Headers          http.Header
	BaseURL          *string
	ProviderMetadata *ProviderMetadata
}

func (o SegmentingOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }
