package api

import "context"

// SparseEmbeddingModel produces sparse embeddings for input texts.
type SparseEmbeddingModel interface {
	SpecificationVersion() string
	ProviderName() string
	ModelID() string
	SupportsParallelCalls() bool
	// SparseEmbed returns a sparse embedding per input text (token->weight map).
	SparseEmbed(ctx context.Context, texts []string, opts SparseEmbeddingOptions) (SparseEmbeddingResponse, error)
}

// SparseEmbeddingResponse contains sparse embeddings for a batch of texts.
type SparseEmbeddingResponse struct {
	RequestID  string
	Embeddings []map[string]float64
}
