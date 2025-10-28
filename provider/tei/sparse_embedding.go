package tei

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/tei/internal/codec"
)

type SparseEmbeddingModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.EmbeddingModel[string, api.SparseEmbedding] = &SparseEmbeddingModel{}

// SparseTextEmbeddingModel creates a new TEI sparse embedding model.
func (p *Provider) SparseTextEmbeddingModel(modelID string) (*SparseEmbeddingModel, error) {
	model := &SparseEmbeddingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.sparse-embedding", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}

	return model, nil
}

func (m *SparseEmbeddingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *SparseEmbeddingModel) SpecificationVersion() string {
	return "v1"
}

func (m *SparseEmbeddingModel) ModelID() string {
	return m.modelID
}

// SupportsParallelCalls returns whether the model can handle multiple embedding calls in parallel.
func (m *SparseEmbeddingModel) SupportsParallelCalls() bool {
	return true
}

// MaxEmbeddingsPerCall returns the limit of how many embeddings can be generated in a single API call.
func (m *SparseEmbeddingModel) MaxEmbeddingsPerCall() *int {
	max := 1000 // TODO: verify
	return &max
}

// DoEmbed generates a list of sparse embeddings for the given input values.
func (m *SparseEmbeddingModel) DoEmbed(
	ctx context.Context,
	values []string,
	opts api.TransportOptions,
) (api.SparseEmbeddingResponse, error) {
	sparseEmbeddingParams, teiOpts, _, err := codec.EncodeSparseEmbedding(
		m.modelID,
		values,
		opts,
	)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}

	resp, err := m.pc.client.Embedding.NewSparse(ctx, sparseEmbeddingParams, teiOpts...)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}

	return codec.DecodeSparseEmbedding(resp)
}
