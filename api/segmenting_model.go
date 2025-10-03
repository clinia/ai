package api

import "context"

// Segment represents a semantically meaningful piece of text.
// For now, mirrors Chunk fields to keep compatibility with existing providers.
type Segment struct {
	ID         string
	Text       string
	StartIndex int
	EndIndex   int
	TokenCount int
}

// SegmentingResponse contains segmentation results for a batch of texts.
// Mirrors ChunkingResponse for initial compatibility.
type SegmentingResponse struct {
	RequestID string
	Segments  [][]Segment
}

// SegmentingModel segments input texts into smaller parts.
// This is analogous to ChunkingModel but adopts the "segment" terminology.
type SegmentingModel interface {
	SpecificationVersion() string
	ProviderName() string
	ModelID() string
	SupportsParallelCalls() bool
	DoSegment(ctx context.Context, texts []string, opts SegmentingOptions) (SegmentingResponse, error)
}
