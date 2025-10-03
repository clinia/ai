package codec

import (
	"fmt"

	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/client/option"
)

// EncodeSegment prepares a Jina Segmenter request for a single text.
// Note: Jina Segmenter segments one content per request; batching is handled by the model.
func EncodeSegment(text string, opts api.SegmentingOptions) (jina.SegmentRequest, []option.RequestOption, error) {
	if text == "" {
		return jina.SegmentRequest{}, nil, fmt.Errorf("jina/segment: content is empty")
	}
	body := jina.SegmentRequest{
		Content:      text,
		ReturnChunks: true,
	}
	var o []option.RequestOption
	if opts.Headers != nil {
		for k, vals := range opts.Headers {
			for _, v := range vals {
				o = append(o, option.WithHeaderAdd(k, v))
			}
		}
	}
	if opts.BaseURL != nil {
		o = append(o, option.WithBaseURL(*opts.BaseURL))
	}
	return body, o, nil
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
