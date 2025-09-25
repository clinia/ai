package clinia

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
)

// ChunkingModel represents a Clinia chunking model.
type ChunkingModel struct {
	modelID      string
	modelName    string
	modelVersion string
	config       ProviderConfig
}

var _ api.ChunkingModel = (*ChunkingModel)(nil)

// ChunkingModel constructs a chunking model wrapper from a model ID ("name:version").
func (p *Provider) ChunkingModel(modelID string) (*ChunkingModel, error) {
	if p.chunker == nil {
		return nil, fmt.Errorf("%s: provider chunker is nil", p.name)
	}

	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &ChunkingModel{
		modelID:      joinModelID(name, version),
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName: p.providerNameFor("chunker"),
			chunker:      p.chunker,
		},
	}, nil
}

func (m *ChunkingModel) SpecificationVersion() string { return "v1" }

func (m *ChunkingModel) ProviderName() string { return m.config.providerName }

func (m *ChunkingModel) ModelID() string { return m.modelID }

func (m *ChunkingModel) SupportsParallelCalls() bool { return true }

func (m *ChunkingModel) Chunk(ctx context.Context, texts []string, opts api.ChunkingOptions) (api.ChunkingResponse, error) {
	params, err := codec.EncodeChunk(texts)
	if err != nil {
		return api.ChunkingResponse{}, err
	}

	if m.config.chunker == nil {
		return api.ChunkingResponse{}, fmt.Errorf("%s: chunker is nil", m.config.providerName)
	}

	resp, err := m.config.chunker.Chunk(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.ChunkingResponse{}, err
	}

	return codec.DecodeChunk(resp)
}
