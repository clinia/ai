package clinia

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
)

// SegmentingModel represents a Clinia segmenting model (previously chunking).
type SegmentingModel struct {
	modelID      string
	modelName    string
	modelVersion string
	config       ProviderConfig
}

var _ api.SegmentingModel = (*SegmentingModel)(nil)

// Segmenter constructs a segmenting model wrapper from a model ID ("name:version").
func (p *Provider) Segmenter(modelID string) (*SegmentingModel, error) {
	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &SegmentingModel{
		modelID:      modelID,
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName:  p.providerNameFor("segmenter"),
			clientOptions: p.clientOptions,
			newChunker:    p.newChunker, // reuse chunker implementation
		},
	}, nil
}

func (m *SegmentingModel) SpecificationVersion() string { return "v1" }
func (m *SegmentingModel) ProviderName() string         { return m.config.providerName }
func (m *SegmentingModel) ModelID() string              { return m.modelID }
func (m *SegmentingModel) SupportsParallelCalls() bool  { return true }

func (m *SegmentingModel) DoSegment(ctx context.Context, texts []string, opts api.SegmentingOptions) (resp api.SegmentingResponse, err error) {
	params, err := codec.EncodeSegment(texts, opts)
	if err != nil {
		return api.SegmentingResponse{}, err
	}
	requester, err := makeRequester(ctx, opts.BaseURL)
	if err != nil {
		return api.SegmentingResponse{}, err
	}
	if requester == nil {
		return api.SegmentingResponse{}, fmt.Errorf("%s: requester is nil", m.config.providerName)
	}

	defer func() {
		if cerr := requester.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if m.config.newChunker == nil {
		return api.SegmentingResponse{}, fmt.Errorf("%s: chunker factory is nil", m.config.providerName)
	}
	chunker := m.config.newChunker(ctx, m.config.clientOptionsWith(requester))
	if chunker == nil {
		return api.SegmentingResponse{}, fmt.Errorf("%s: chunker factory returned nil", m.config.providerName)
	}

	respProto, err := chunker.Chunk(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.SegmentingResponse{}, err
	}
	return codec.DecodeSegment(respProto)
}
