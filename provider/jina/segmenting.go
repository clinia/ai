package jina

import (
	"context"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/jina/internal/codec"
)

// SegmenterModel implements api.SegmentingModel using Jina Segmenter API.
type SegmenterModel struct {
	modelID string
	pc      ProviderConfig
}

var _ api.SegmentingModel = &SegmenterModel{}

// Segmenter creates a new SegmenterModel.
func (p *Provider) Segmenter(modelID string) (api.SegmentingModel, error) {
	m := &SegmenterModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: p.name + ".segmenter",
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}
	return m, nil
}

func (m *SegmenterModel) SpecificationVersion() string { return "v1" }
func (m *SegmenterModel) ProviderName() string         { return m.pc.providerName }
func (m *SegmenterModel) ModelID() string              { return m.modelID }
func (m *SegmenterModel) SupportsParallelCalls() bool  { return true }

func (m *SegmenterModel) DoSegment(ctx context.Context, texts []string, opts api.SegmentingOptions) (api.SegmentingResponse, error) {
	// Inspect provider metadata for batching preference.
	if meta := codec.GetSegmenterMetadata(opts); meta != nil && meta.UseContentArray {
		// True batching: send all texts in one request using content as []string
		body, ropts, err := codec.EncodeSegmentBatch(texts, opts)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		resp, err := m.pc.client.Segmenter.NewBatch(ctx, body, ropts...)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		groups, err := codec.DecodeSegmentBatch(resp)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		return api.SegmentingResponse{Segments: groups}, nil
	}

	// Default (and Jina-official) behavior: one content per request; iterate.
	groups := make([][]api.Segment, 0, len(texts))
	for _, text := range texts {
		body, ropts, err := codec.EncodeSegment(text, opts)
		if err != nil {
			return api.SegmentingResponse{}, err
		}
		resp, err := m.pc.client.Segmenter.New(ctx, body, ropts...)
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
