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

var _ api.EmbeddingModel[string, api.SparseEmbedding] = (*SparseEmbeddingModel)(nil)

func (p *Provider) SparseEmbeddingModel(modelID string) (*SparseEmbeddingModel, error) {
	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &SparseEmbeddingModel{
		modelID:      modelID,
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName:  p.providerNameFor("sparse_embedding"),
			clientOptions: p.clientOptions,
			newSparse:     p.newSparse,
		},
	}, nil
}

func (m *SparseEmbeddingModel) SpecificationVersion() string { return "v1" }
func (m *SparseEmbeddingModel) ProviderName() string         { return m.config.providerName }
func (m *SparseEmbeddingModel) ModelID() string              { return m.modelID }
func (m *SparseEmbeddingModel) SupportsParallelCalls() bool  { return true }
func (m *SparseEmbeddingModel) MaxEmbeddingsPerCall() *int   { return nil }

func (m *SparseEmbeddingModel) DoEmbed(ctx context.Context, texts []string, opts api.TransportOptions) (resp api.SparseEmbeddingResponse, err error) {
	params, err := codec.EncodeSparseEmbedding(texts, opts)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}

	requester, err := makeRequester(ctx, opts.BaseURL)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}
	if requester == nil {
		return api.SparseEmbeddingResponse{}, fmt.Errorf("%s: requester is nil", m.config.providerName)
	}

	defer func() {
		if cerr := requester.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if m.config.newSparse == nil {
		return api.SparseEmbeddingResponse{}, fmt.Errorf("%s: sparse embedder factory is nil", m.config.providerName)
	}

	sparse := m.config.newSparse(ctx, m.config.clientOptionsWith(requester))
	if sparse == nil {
		return api.SparseEmbeddingResponse{}, fmt.Errorf("%s: sparse embedder factory returned nil", m.config.providerName)
	}

	res, err := sparse.SparseEmbed(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}

	resp, err = codec.DecodeSparseEmbedding(res)
	if err != nil {
		return api.SparseEmbeddingResponse{}, err
	}
	return resp, nil
}
