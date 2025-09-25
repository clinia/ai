package clinia

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
)

// SparseEmbeddingModel wraps the Clinia sparse embedder.
type SparseEmbeddingModel struct {
	modelID      string
	modelName    string
	modelVersion string
	config       ProviderConfig
}

var _ api.SparseEmbeddingModel = (*SparseEmbeddingModel)(nil)

func (p *Provider) SparseEmbeddingModel(modelID string) (*SparseEmbeddingModel, error) {
	if p.sparse == nil {
		return nil, fmt.Errorf("%s: provider sparse embedder is nil", p.name)
	}

	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &SparseEmbeddingModel{
		modelID:      joinModelID(name, version),
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName: p.providerNameFor("sparse_embedding"),
			sparse:       p.sparse,
		},
	}, nil
}

func (m *SparseEmbeddingModel) SpecificationVersion() string { return "v1" }
func (m *SparseEmbeddingModel) ProviderName() string         { return m.config.providerName }
func (m *SparseEmbeddingModel) ModelID() string              { return m.modelID }
func (m *SparseEmbeddingModel) SupportsParallelCalls() bool  { return true }

func (m *SparseEmbeddingModel) SparseEmbed(ctx context.Context, texts []string, opts api.SparseEmbeddingOptions) (api.SparseEmbeddingResponse, error) {
	params, err := codec.EncodeSparseEmbedding(m.modelName, m.modelVersion, texts, opts)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}
	if m.config.sparse == nil {
		return api.SparseEmbeddingResponse{}, fmt.Errorf("%s: sparse embedder is nil", m.config.providerName)
	}
	res, err := m.config.sparse.SparseEmbed(ctx, params.ModelName, params.ModelVersion, params.Request)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}
	return codec.DecodeSparseEmbedding(res)
}
