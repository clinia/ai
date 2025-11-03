package chonkie

import (
	"context"

	"go.jetify.com/ai/provider/internal/requesterx"
)

// SegmentRequest models the POST body for Chonkie Segmenting API.
type SegmentRequest struct {
	Content string `json:"content"`
	segmentCommon
}

// BatchSegmentRequest mirrors SegmentRequest but allows sending multiple
// contents in a single request by using an array for the content field.
type BatchSegmentRequest struct {
	Content []string `json:"content"`
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
// segmenting behavior for Chonkie. Currently supports enabling true batching
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

// New performs a POST /segment request with a batched body where
// content is an array of strings. The response is expected to be an array
// of SegmentResponse objects, one per input string.
func (s SegmentingService) New(ctx context.Context, body BatchSegmentRequest, opts ...requesterx.RequestOption) (res []SegmentResponse, err error) {
	all := append([]requesterx.RequestOption{}, s.opts...)
	all = append(all, opts...)
	err = requesterx.ExecuteNewRequest(ctx, "POST", "", body, &res, all...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
