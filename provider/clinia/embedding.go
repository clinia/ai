package clinia

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
)

// EmbeddingModel represents a Clinia embedding model backed by the models-client-go embedder.
type EmbeddingModel struct {
	modelID      string
	modelName    string
	modelVersion string
	config       ProviderConfig
}

var _ api.EmbeddingModel[string, api.Embedding] = (*EmbeddingModel)(nil)

// TextEmbeddingModel constructs a new text embedding model from a model ID in the form "name:version".
// Implements api.Provider.TextEmbeddingModel.
func (p *Provider) TextEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.Embedding], error) {
	name, version, err := splitModelID("embed", modelID)
	if err != nil {
		return nil, err
	}

	model := &EmbeddingModel{
		modelID:      modelID,
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName:  p.providerNameFor("embedding"),
			clientOptions: p.clientOptions,
			newEmbedder:   p.newEmbedder,
		},
	}
	return model, nil
}

func (m *EmbeddingModel) ProviderName() string {
	return m.config.providerName
}

func (m *EmbeddingModel) SpecificationVersion() string {
	return "v2"
}

func (m *EmbeddingModel) ModelID() string {
	return m.modelID
}

func (m *EmbeddingModel) SupportsParallelCalls() bool {
	return true
}

func (m *EmbeddingModel) MaxEmbeddingsPerCall() *int {
	return nil
}

// DoEmbed executes an embedding call against the Clinia embedder.
func (m *EmbeddingModel) DoEmbed(ctx context.Context, values []string, opts api.EmbeddingOptions) (resp api.DenseEmbeddingResponse, err error) {
	params, err := codec.EncodeEmbedding(values, opts)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}
	requester, err := makeRequester(ctx, opts.BaseURL)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}
	if requester == nil {
		return api.DenseEmbeddingResponse{}, fmt.Errorf("clinia/embed: requester is nil")
	}

	defer func() {
		if cerr := requester.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if m.config.newEmbedder == nil {
		return api.DenseEmbeddingResponse{}, fmt.Errorf("clinia/embed: embedder factory is nil")
	}

	embedder := m.config.newEmbedder(ctx, m.config.clientOptionsWith(requester))
	if embedder == nil {
		return api.DenseEmbeddingResponse{}, fmt.Errorf("clinia/embed: embedder factory returned nil")
	}

	embedResp, err := embedder.Embed(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}

	resp, err = codec.DecodeEmbedding(embedResp)
	if err != nil {
		return api.DenseEmbeddingResponse{}, err
	}
	return resp, nil
}
