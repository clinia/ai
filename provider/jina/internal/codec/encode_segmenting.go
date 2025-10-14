package codec

import (
	"fmt"

	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/client/option"
)

// EncodeSegment prepares a Jina Segmenter request for a single text.
// Note: Jina Segmenter segments one content per request; batching is handled by the model.
func EncodeSegment(text string, opts api.TransportOptions) (jina.SegmentRequest, []option.RequestOption, error) {
	if text == "" {
		return jina.SegmentRequest{}, nil, fmt.Errorf("jina/segment: content is empty")
	}

	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, option.WithAPIKey(opts.APIKey))
	}

	if opts.BaseURL != nil {
		reqOpts = append(reqOpts, option.WithBaseURL(*opts.BaseURL))
	}

	if opts.UseRawBaseURL {
		reqOpts = append(reqOpts, option.WithUseRawBaseURL())
	}

	body := jina.SegmentRequest{Content: text}
	body.ReturnChunks = true

	return body, reqOpts, nil
}

// EncodeSegmentBatch prepares a Jina Segmenter request for multiple texts
// in a single HTTP call by using an array in the "content" field.
func EncodeSegmentBatch(texts []string, opts api.TransportOptions) (jina.BatchSegmentRequest, []option.RequestOption, error) {
	if len(texts) == 0 {
		return jina.BatchSegmentRequest{}, nil, fmt.Errorf("jina/segment: texts cannot be empty")
	}
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, option.WithAPIKey(opts.APIKey))
	}

	if opts.BaseURL != nil {
		reqOpts = append(reqOpts, option.WithBaseURL(*opts.BaseURL))
	}

	if opts.UseRawBaseURL {
		reqOpts = append(reqOpts, option.WithUseRawBaseURL())
	}

	body := jina.BatchSegmentRequest{Content: texts}
	body.ReturnChunks = true
	return body, reqOpts, nil
}

// DecodeSegment maps a Jina Segmenter response to a list of SDK segments.
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

// DecodeSegmentBatch maps a batched Jina Segmenter response (slice) to the
// [][]api.Segment shape expected by the SDK.
func DecodeSegmentBatch(resps *[]jina.SegmentResponse) ([][]api.Segment, error) {
	if resps == nil || *resps == nil {
		return nil, fmt.Errorf("jina/segment: batch response is nil")
	}
	out := make([][]api.Segment, 0, len(*resps))
	for i := range *resps {
		segs, err := DecodeSegment(&(*resps)[i])
		if err != nil {
			return nil, err
		}
		out = append(out, segs)
	}
	return out, nil
}
