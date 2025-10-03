package jina

import (
	"context"

	"go.jetify.com/ai/provider/jina/client/internal/requestconfig"
	"go.jetify.com/ai/provider/jina/client/option"
)

// SegmentRequest models the POST body for Jina Segmenter API.
type SegmentRequest struct {
	Content        string  `json:"content"`
	ReturnTokens   bool    `json:"return_tokens,omitempty"`
	ReturnChunks   bool    `json:"return_chunks,omitempty"`
	MaxChunkLength *int    `json:"max_chunk_length,omitempty"`
	Head           *int    `json:"head,omitempty"`
	Tail           *int    `json:"tail,omitempty"`
	Tokenizer      *string `json:"tokenizer,omitempty"`
}

// SegmentResponse is a subset of the Segmenter API response needed to build segments.
type SegmentResponse struct {
	NumTokens      int      `json:"num_tokens"`
	Tokenizer      string   `json:"tokenizer"`
	NumChunks      int      `json:"num_chunks"`
	ChunkPositions [][]int  `json:"chunk_positions"`
	Chunks         []string `json:"chunks"`
}

type SegmenterService struct{ opts []option.RequestOption }

func NewSegmenterService(opts ...option.RequestOption) SegmenterService {
	return SegmenterService{opts: opts}
}

// New performs a POST /segment request.
func (s SegmenterService) New(ctx context.Context, body SegmentRequest, opts ...option.RequestOption) (res *SegmentResponse, err error) {
	all := append([]option.RequestOption{}, s.opts...)
	all = append(all, opts...)
	err = requestconfig.ExecuteNewRequest(ctx, "POST", "segment", body, &res, all...)
	return
}
