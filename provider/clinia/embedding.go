package clinia

import (
	"context"
	"fmt"
	"strings"

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

var _ api.EmbeddingModel[string] = (*EmbeddingModel)(nil)

// NewEmbeddingModel constructs a new embedding model wrapper.
func (p *Provider) NewEmbeddingModel(modelName, modelVersion string) (*EmbeddingModel, error) {
	if p.embedder == nil {
		return nil, fmt.Errorf("clinia/embed: provider embedder is nil")
	}

	name := strings.TrimSpace(modelName)
	if name == "" {
		return nil, fmt.Errorf("clinia/embed: model name is required")
	}

	version := strings.TrimSpace(modelVersion)
	if version == "" {
		return nil, fmt.Errorf("clinia/embed: model version is required")
	}

	model := &EmbeddingModel{
		modelID:      buildModelID(name, version),
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName: p.providerNameFor("embedding"),
			embedder:     p.embedder,
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
func (m *EmbeddingModel) DoEmbed(ctx context.Context, values []string, opts api.EmbeddingOptions) (api.EmbeddingResponse, error) {
	params, err := codec.EncodeEmbedding(m.modelName, m.modelVersion, values, opts)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	resp, err := m.config.embedder.Embed(ctx, params.ModelName, params.ModelVersion, params.Request)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	return codec.DecodeEmbedding(resp)
}
