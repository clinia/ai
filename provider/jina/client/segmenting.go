package jina

import (
	"context"

	"go.jetify.com/ai/provider/internal/requesterx"
)

// SegmentRequest models the POST body for Jina Segmenting API.
type SegmentRequest struct {
	Content string `json:"content"`
	segmentCommon
}

// segmentCommon holds attributes shared by both single and batched requests.
// Embedded into SegmentRequest and BatchSegmentRequest to avoid duplication.
type segmentCommon struct {
	ReturnTokens   bool    `json:"return_tokens,omitempty"`
	ReturnChunks   bool    `json:"return_chunks,omitempty"`
	MaxChunkLength *int    `json:"max_chunk_length,omitempty"`
	Head           *int    `json:"head,omitempty"`
	Tail           *int    `json:"tail,omitempty"`
	Tokenizer      *string `json:"tokenizer,omitempty"`
}

// SegmentingNewParams allows callers to pass provider metadata to tweak
// segmenting behavior for Jina. Currently supports enabling true batching
// by sending the content as an array in a single request.
type SegmentingNewParams struct {
	// UseContentArray toggles sending a batched request with content as []string.
	// If false, provider will send one request per input.
	UseContentArray bool `json:"use_content_array,omitempty"`
}

// SegmentResponse is a subset of the Segmenting API response needed to build segments.
type SegmentResponse struct {
	NumTokens      int      `json:"num_tokens"`
	Tokenizer      string   `json:"tokenizer"`
	NumChunks      int      `json:"num_chunks"`
	ChunkPositions [][]int  `json:"chunk_positions"`
	Chunks         []string `json:"chunks"`
}

type SegmentingService struct{ opts []requesterx.RequestOption }

func NewSegmentingService(opts ...requesterx.RequestOption) SegmentingService {
	return SegmentingService{opts: opts}
}

// New performs a POST /segment request.
func (s SegmentingService) New(ctx context.Context, body SegmentRequest, opts ...requesterx.RequestOption) (res *SegmentResponse, err error) {
	all := append([]requesterx.RequestOption{}, s.opts...)
	all = append(all, opts...)
	path := "segment"
	err = requesterx.ExecuteNewRequest(ctx, "POST", path, body, &res, all...)
	return res, err
}
