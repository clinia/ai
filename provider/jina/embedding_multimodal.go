package jina

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/internal/codec"
)

// EmbeddingModel represents an Jina embedding model.
type MultimodalEmbeddingModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.EmbeddingModel[jina.MultimodalEmbeddingInput] = &MultimodalEmbeddingModel{}

// NewEmbeddingModel creates a new Jina embedding model.
func (p *Provider) NewMultimodalEmbeddingModel(modelID string) *MultimodalEmbeddingModel {
	// Create model with provider's client
	model := &MultimodalEmbeddingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.embedding", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}

	return model
}

func (m *MultimodalEmbeddingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *MultimodalEmbeddingModel) SpecificationVersion() string {
	return "v2"
}

func (m *MultimodalEmbeddingModel) ModelID() string {
	return m.modelID
}

// SupportsParallelCalls implements api.EmbeddingModel.
func (m *MultimodalEmbeddingModel) SupportsParallelCalls() bool {
	return true
}

// MaxEmbeddingsPerCall implements api.EmbeddingModel.
func (m *MultimodalEmbeddingModel) MaxEmbeddingsPerCall() *int {
	max := 32768
	return &max
}

// DoEmbed implements api.EmbeddingModel.
func (m *MultimodalEmbeddingModel) DoEmbed(
	ctx context.Context,
	values []jina.MultimodalEmbeddingInput,
	opts api.EmbeddingOptions,
) (api.EmbeddingResponse, error) {
	embeddingParams, jinaOpts, _, err := codec.EncodeMultimodalEmbedding(
		m.modelID,
		values,
		opts,
	)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	resp, err := m.pc.client.Embeddings.NewMultiModal(ctx, embeddingParams, jinaOpts...)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	return codec.DecodeEmbedding(resp)
}
