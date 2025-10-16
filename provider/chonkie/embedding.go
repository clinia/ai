package chonkie

import (
	"go.jetify.com/ai/api"
)

// TextEmbeddingModel creates a new Chonkie embedding model.
func (p *Provider) TextEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "TextEmbeddingModel")
}
