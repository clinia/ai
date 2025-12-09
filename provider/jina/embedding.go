package jina

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/instrumentation"
	"go.jetify.com/ai/provider/jina/internal/codec"
)

// EmbeddingModel represents an Jina embedding model.
type EmbeddingModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.EmbeddingModel[string, api.Embedding] = &EmbeddingModel{}

// TextEmbeddingModel creates a new Jina embedding model.
func (p *Provider) TextEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.Embedding], error) {
	// Create model with provider's client
	model := &EmbeddingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.embedding", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
			instrumenter: p.instrumenter,
		},
	}

	return model, nil
}

func (m *EmbeddingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *EmbeddingModel) SpecificationVersion() string {
	return "v2"
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
	max := 32768
	return &max
}

// DoEmbed implements api.EmbeddingModel.
func (m *EmbeddingModel) DoEmbed(
	ctx context.Context,
	values []string,
	opts api.TransportOptions,
) (resp api.DenseEmbeddingResponse, err error) {
	ctx, span := m.pc.instrumenter.Start(
		ctx,
		"DoEmbed",
		instrumentation.Attributes{
			"provider":   m.ProviderName(),
			"model":      m.modelID,
			"model_type": "embedding",
			"operation":  string(instrumentation.OperationEmbed),
		},
		instrumentation.ProviderSpanInfo{
			Provider:  m.ProviderName(),
			Model:     m.modelID,
			Operation: instrumentation.OperationEmbed,
		},
	)
	defer instrumentation.EndSpan(span, &err)

	embeddingParams, jinaOpts, _, err := codec.EncodeEmbedding(
		m.modelID,
		values,
		opts,
	)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}

	apiResp, err := m.pc.client.Embeddings.New(ctx, embeddingParams, jinaOpts...)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}

	return codec.DecodeEmbedding(apiResp)
}
