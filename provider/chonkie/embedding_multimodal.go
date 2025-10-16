package chonkie

import (
	"go.jetify.com/ai/api"
)

// NewEmbeddingModel creates a new Chonkie embedding model.
func (p *Provider) MultimodalEmbeddingModel(modelID string) (api.EmbeddingModel[api.MultimodalEmbeddingInput, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "MultimodalEmbeddingModel")
}
