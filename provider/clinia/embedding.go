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

// TextEmbeddingModel constructs a new text embedding model from a model ID in the form "name:version".
// Implements api.Provider.TextEmbeddingModel.
func (p *Provider) TextEmbeddingModel(modelID string) (api.EmbeddingModel[string], error) {
    if p.embedder == nil {
        return nil, fmt.Errorf("clinia/embed: provider embedder is nil")
    }

    id := strings.TrimSpace(modelID)
    if id == "" {
        return nil, fmt.Errorf("clinia/embed: model id is required (expected 'name:version')")
    }

    // Require explicit version; split on first ':'
    parts := strings.SplitN(id, ":", 2)
    if len(parts) != 2 {
        return nil, fmt.Errorf("clinia/embed: model version is required in id (expected 'name:version')")
    }

    name := strings.TrimSpace(parts[0])
    version := strings.TrimSpace(parts[1])
    if name == "" {
        return nil, fmt.Errorf("clinia/embed: model name is required")
    }
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
	params, err := codec.EncodeEmbedding(values, opts)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	resp, err := m.config.embedder.Embed(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.EmbeddingResponse{}, err
	}

	return codec.DecodeEmbedding(resp)
}
