package clinia

import (
	"context"
	"fmt"
	"strings"

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

func (p *Provider) NewSparseEmbeddingModel(modelName, modelVersion string) (*SparseEmbeddingModel, error) {
	if p.sparse == nil {
		return nil, fmt.Errorf("clinia/sparse: provider sparse embedder is nil")
	}

	name := strings.TrimSpace(modelName)
	if name == "" {
		return nil, fmt.Errorf("clinia/sparse: model name is required")
	}

	version := strings.TrimSpace(modelVersion)
	if version == "" {
		return nil, fmt.Errorf("clinia/sparse: model version is required")
	}

	return &SparseEmbeddingModel{
		modelID:      buildModelID(modelName, modelVersion),
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
		return api.SparseEmbeddingResponse{}, fmt.Errorf("clinia/sparse: sparse embedder is nil")
	}
	res, err := m.config.sparse.SparseEmbed(ctx, params.ModelName, params.ModelVersion, params.Request)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}
	return codec.DecodeSparseEmbedding(res)
}
