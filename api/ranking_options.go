package api

import "net/http"

// RankingOptions represents per-call options for ranking models.
type RankingOptions struct {
	Headers          http.Header
	BaseURL          *string
	ProviderMetadata *ProviderMetadata
}

func (o RankingOptions) GetProviderMetadata() *ProviderMetadata { return o.ProviderMetadata }
