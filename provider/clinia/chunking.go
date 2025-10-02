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
	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &ChunkingModel{
		modelID:      modelID,
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName:  p.providerNameFor("chunker"),
			clientOptions: p.clientOptions,
			newChunker:    p.newChunker,
		},
	}, nil
}

func (m *ChunkingModel) SpecificationVersion() string { return "v1" }

func (m *ChunkingModel) ProviderName() string { return m.config.providerName }

func (m *ChunkingModel) ModelID() string { return m.modelID }

func (m *ChunkingModel) SupportsParallelCalls() bool { return true }

func (m *ChunkingModel) DoChunk(ctx context.Context, texts []string, opts api.ChunkingOptions) (resp api.ChunkingResponse, err error) {
	params, err := codec.EncodeChunk(texts, opts)
	if err != nil {
		return api.ChunkingResponse{}, err
	}

	requester, err := makeRequester(ctx, opts.BaseURL)
	if err != nil {
		return api.ChunkingResponse{}, err
	}
	if requester == nil {
		return api.ChunkingResponse{}, fmt.Errorf("%s: requester is nil", m.config.providerName)
	}

	defer func() {
		if cerr := requester.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if m.config.newChunker == nil {
		return api.ChunkingResponse{}, fmt.Errorf("%s: chunker factory is nil", m.config.providerName)
	}

	chunker := m.config.newChunker(ctx, m.config.clientOptionsWith(requester))
	if chunker == nil {
		return api.ChunkingResponse{}, fmt.Errorf("%s: chunker factory returned nil", m.config.providerName)
	}

	respProto, err := chunker.Chunk(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.ChunkingResponse{}, err
	}

	resp, err = codec.DecodeChunk(respProto)
	if err != nil {
		return api.ChunkingResponse{}, err
	}
	return resp, nil
}
