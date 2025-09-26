package clinia

import (
	"github.com/clinia/models-client-go/cliniamodel/common"
	"go.jetify.com/ai"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
)

// WithEmbeddingRequester injects a pre-built requester for an embedding call.
func WithEmbeddingRequester(r common.Requester) ai.EmbeddingOption {
	return func(o *api.EmbeddingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(nil)
		}
		o.ProviderMetadata.Set("clinia", codec.Metadata{Requester: r})
	}
}

// WithRankingRequester injects a pre-built requester for a ranking call.
func WithRankingRequester(r common.Requester) ai.RankingOption {
	return func(o *api.RankingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(nil)
		}
		o.ProviderMetadata.Set("clinia", codec.Metadata{Requester: r})
	}
}

// WithChunkingRequester injects a pre-built requester for a chunking call.
func WithChunkingRequester(r common.Requester) ai.ChunkingOption {
	return func(o *api.ChunkingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(nil)
		}
		o.ProviderMetadata.Set("clinia", codec.Metadata{Requester: r})
	}
}

// WithSparseEmbeddingRequester injects a pre-built requester for a sparse embedding call.
func WithSparseEmbeddingRequester(r common.Requester) ai.SparseEmbeddingOption {
	return func(o *api.SparseEmbeddingOptions) {
		if o.ProviderMetadata == nil {
			o.ProviderMetadata = api.NewProviderMetadata(nil)
		}
		o.ProviderMetadata.Set("clinia", codec.Metadata{Requester: r})
	}
}
