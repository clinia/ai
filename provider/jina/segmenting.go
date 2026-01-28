package jina

import (
	"context"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/instrumentation"
	"go.jetify.com/ai/provider/jina/internal/codec"
)

// SegmentingModel implements api.SegmentingModel using Jina Segmenting API.
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
			instrumenter: p.instrumenter,
		},
	}
	return m, nil
}

func (m *SegmentingModel) SpecificationVersion() string { return "v1" }
func (m *SegmentingModel) ProviderName() string         { return m.pc.providerName }
func (m *SegmentingModel) ModelID() string              { return m.modelID }
func (m *SegmentingModel) SupportsParallelCalls() bool  { return true }

func (m *SegmentingModel) DoSegment(ctx context.Context, texts []string, opts api.TransportOptions) (resp api.SegmentingResponse, err error) {
	ctx, span := m.pc.instrumenter.Start(
		ctx,
		"DoSegment",
		instrumentation.Attributes{
			"provider":   m.ProviderName(),
			"model":      m.modelID,
			"model_type": "segmenting",
			"operation":  string(instrumentation.OperationSegment),
		},
		instrumentation.ProviderSpanInfo{
			Provider:  m.ProviderName(),
			Model:     m.modelID,
			Operation: instrumentation.OperationSegment,
		},
	)
	defer instrumentation.EndSpan(span, &err)

	groups := make([][]api.Segment, 0, len(texts))
	for _, text := range texts {
		body, ropts, err := codec.EncodeSegment(text, opts)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		resp, err := m.pc.client.Segments.New(ctx, body, ropts...)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		segs, err := codec.DecodeSegment(resp)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		groups = append(groups, segs)
	}
	return api.SegmentingResponse{Segments: groups}, nil
}
