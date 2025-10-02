package api

import "context"

// Chunk represents a single chunk produced by a chunking model.
type Chunk struct {
	ID         string
	Text       string
	StartIndex int
	EndIndex   int
	TokenCount int
}

// ChunkingResponse contains chunking results for a batch of texts.
type ChunkingResponse struct {
	// RequestID mirrors the ID supplied in the request when supported.
	RequestID string
	// Chunks is grouped per input text (outer slice aligns with input order).
	Chunks [][]Chunk
}

// ChunkingModel is a model that splits texts into chunks.
type ChunkingModel interface {
	SpecificationVersion() string
	ProviderName() string
	ModelID() string
	SupportsParallelCalls() bool
	DoChunk(ctx context.Context, texts []string, opts ChunkingOptions) (ChunkingResponse, error)
}
