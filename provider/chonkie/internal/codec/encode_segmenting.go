package codec

import (
	"fmt"

	"go.jetify.com/ai/api"
	chonkie "go.jetify.com/ai/provider/chonkie/client"
	"go.jetify.com/ai/provider/internal/requesterx"
)

// EncodeSegment prepares a Chonkie Segmenting request for a single text.
// Note: Chonkie Segmenting segments one content per request; batching is handled by the model.
func EncodeSegment(text string, opts api.TransportOptions) (chonkie.SegmentRequest, []requesterx.RequestOption, error) {
	if text == "" {
		return chonkie.SegmentRequest{}, nil, fmt.Errorf("chonkie/segment: content is empty")
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

	body := chonkie.SegmentRequest{Content: text}
	body.ReturnChunks = true

	return body, reqOpts, nil
}

// EncodeSegmentBatch prepares a Chonkie Segmenting request for multiple texts
// in a single HTTP call by using an array in the "content" field.
func EncodeSegmentBatch(texts []string, opts api.TransportOptions) (chonkie.BatchSegmentRequest, []requesterx.RequestOption, error) {
	if len(texts) == 0 {
		return chonkie.BatchSegmentRequest{}, nil, fmt.Errorf("chonkie/segment: texts cannot be empty")
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

	body := chonkie.BatchSegmentRequest{Content: texts}
	body.ReturnChunks = true
	return body, reqOpts, nil
}

// DecodeSegment maps a Chonkie Segmenting response to a list of SDK segments.
func DecodeSegment(resp *chonkie.SegmentResponse) ([]api.Segment, error) {
	if resp == nil {
		return nil, fmt.Errorf("chonkie/segment: response is nil")
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

// DecodeSegmentBatch maps a batched Chonkie Segmenting response (slice) to the
// [][]api.Segment shape expected by the SDK.
func DecodeSegmentBatch(resps []chonkie.SegmentResponse) ([][]api.Segment, error) {
	if len(resps) == 0 {
		return nil, fmt.Errorf("chonkie/segment: batch response is nil")
	}
	out := make([][]api.Segment, 0, len(resps))
	for i := range resps {
		segs, err := DecodeSegment(&(resps)[i])
		if err != nil {
			return nil, err
		}
		out = append(out, segs)
	}
	return out, nil
}
