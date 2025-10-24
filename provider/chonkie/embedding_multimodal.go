package chonkie

import (
	"go.jetify.com/ai/api"
)

// MultimodalEmbeddingModel returns an error as Chonkie does not support multimodal embedding models.
func (p *Provider) MultimodalEmbeddingModel(modelID string) (api.EmbeddingModel[api.MultimodalEmbeddingInput, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "MultimodalEmbeddingModel")
}
