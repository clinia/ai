package chonkie

import (
	"context"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/chonkie/internal/codec"
)

// SegmentingModel implements api.SegmentingModel using Chonkie Segmenting API.
type SegmentingModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.SegmentingModel = &SegmentingModel{}

// SegmentingModel creates a new SegmentingModel.
func (p *Provider) SegmentingModel(modelID string) (api.SegmentingModel, error) {
	m := &SegmentingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: p.name + ".segmenting",
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}
	return m, nil
}

func (m *SegmentingModel) SpecificationVersion() string { return "v1" }
func (m *SegmentingModel) ProviderName() string         { return m.pc.providerName }
func (m *SegmentingModel) ModelID() string              { return m.modelID }
func (m *SegmentingModel) SupportsParallelCalls() bool  { return true }

func (m *SegmentingModel) DoSegment(ctx context.Context, texts []string, opts api.TransportOptions) (api.SegmentingResponse, error) {
	body, ropts, err := codec.EncodeSegmentBatch(texts, opts)
	if err != nil {
		return api.SegmentingResponse{}, err
	}
	resp, err := m.pc.client.Segments.New(ctx, body, ropts...)
	if err != nil {
		return api.SegmentingResponse{}, err
	}
	groups, err := codec.DecodeSegmentBatch(resp)
	if err != nil {
		return api.SegmentingResponse{}, err
	}
	return api.SegmentingResponse{Segments: groups}, nil
}
