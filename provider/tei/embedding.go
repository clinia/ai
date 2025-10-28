package tei

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/tei/internal/codec"
)

// EmbeddingModel represents a TEI embedding model.
type EmbeddingModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.EmbeddingModel[string, api.Embedding] = &EmbeddingModel{}

// TextEmbeddingModel creates a new TEI embedding model.
func (p *Provider) TextEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.Embedding], error) {
	// Create model with provider's client
	model := &EmbeddingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.embedding", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}

	return model, nil
}

func (m *EmbeddingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *EmbeddingModel) SpecificationVersion() string {
	return "v1"
}

func (m *EmbeddingModel) ModelID() string {
	return m.modelID
}

// SupportsParallelCalls implements api.EmbeddingModel.
func (m *EmbeddingModel) SupportsParallelCalls() bool {
	return true
}

// MaxEmbeddingsPerCall implements api.EmbeddingModel.
func (m *EmbeddingModel) MaxEmbeddingsPerCall() *int {
	max := 1000
	return &max
}

// DoEmbed implements api.EmbeddingModel.
func (m *EmbeddingModel) DoEmbed(
	ctx context.Context,
	values []string,
	opts api.TransportOptions,
) (api.DenseEmbeddingResponse, error) {
	embeddingParams, teiOpts, _, err := codec.EncodeEmbedding(
		m.modelID,
		values,
		opts,
	)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}

	resp, err := m.pc.client.Embedding.New(ctx, embeddingParams, teiOpts...)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}

	return codec.DecodeEmbedding(resp)
}
