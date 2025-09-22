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
func (p *Provider) NewEmbeddingModel(modelID string) (*EmbeddingModel, error) {
	const defaultModelVersion = "1"

	if p.embedder == nil {
		return nil, fmt.Errorf("clinia/embed: provider embedder is nil")
	}

	name, version := splitModelIdentifier(modelID, defaultModelVersion)

	model := &EmbeddingModel{
		modelID:      modelID,
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

func splitModelIdentifier(modelID string, defaultVersion string) (string, string) {
	trimmed := strings.TrimSpace(modelID)
	if trimmed == "" {
		return "", defaultVersion
	}

	parts := strings.Split(trimmed, ":")
	if len(parts) > 1 {
		name := strings.TrimSpace(parts[0])
		version := strings.TrimSpace(parts[1])
		if name == "" {
			name = trimmed
		}
		if version == "" {
			version = defaultVersion
		}
		return name, version
	}

	return trimmed, defaultVersion
}
