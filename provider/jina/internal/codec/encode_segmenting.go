package codec

import (
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/requesterx"
	jina "go.jetify.com/ai/provider/jina/client"
)

// EncodeSegment prepares a Jina Segmenting request for a single text.
func EncodeSegment(text string, opts api.TransportOptions) (jina.SegmentRequest, []requesterx.RequestOption, error) {
	if text == "" {
		return jina.SegmentRequest{}, nil, fmt.Errorf("jina/segment: content is empty")
	}

	var reqOpts []requesterx.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, requesterx.WithAPIKey(opts.APIKey))
	}

	if len(opts.BaseURL) > 0 {
		reqOpts = append(reqOpts, requesterx.WithBaseURL(opts.BaseURL))
	}

	if opts.UseRawBaseURL {
		reqOpts = append(reqOpts, requesterx.WithUseRawBaseURL())
	}

	body := jina.SegmentRequest{Content: text}
	body.ReturnChunks = true

	return body, reqOpts, nil
}

// DecodeSegment maps a Jina Segmenting response to a list of SDK segments.
func DecodeSegment(resp *jina.SegmentResponse) ([]api.Segment, error) {
	if resp == nil {
		return nil, fmt.Errorf("jina/segment: response is nil")
	}
	segs := make([]api.Segment, 0, len(resp.Chunks))
	for i, text := range resp.Chunks {
		seg := api.Segment{ID: fmt.Sprintf("c%d", i), Text: text}
		if i < len(resp.ChunkPositions) && len(resp.ChunkPositions[i]) == 2 {
			seg.StartIndex = resp.ChunkPositions[i][0]
			seg.EndIndex = resp.ChunkPositions[i][1]
		}
		segs = append(segs, seg)
	}
	return segs, nil
}
