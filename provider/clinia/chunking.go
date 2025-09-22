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

// NewChunkingModel constructs a chunking model wrapper.
func (p *Provider) NewChunkingModel(modelID string) (*ChunkingModel, error) {
	const defaultModelVersion = "1"

	if p.chunker == nil {
		return nil, fmt.Errorf("clinia/chunk: provider chunker is nil")
	}

	name, version := splitModelIdentifier(modelID, defaultModelVersion)

	return &ChunkingModel{
		modelID:      modelID,
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
	params, err := codec.EncodeChunk(m.modelName, m.modelVersion, texts)
	if err != nil {
		return api.ChunkingResponse{}, err
	}

	if m.config.chunker == nil {
		return api.ChunkingResponse{}, fmt.Errorf("clinia/chunk: chunker is nil")
	}

	resp, err := m.config.chunker.Chunk(ctx, params.ModelName, params.ModelVersion, params.Request)
	if err != nil {
		return api.ChunkingResponse{}, err
	}

	return codec.DecodeChunk(resp)
}
