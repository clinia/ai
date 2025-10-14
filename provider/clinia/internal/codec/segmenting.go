package codec

import (
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"go.jetify.com/ai/api"
)

// SegmentParams holds the Clinia segment request (reuses chunk request DTO).
type SegmentParams struct {
	Request cliniaclient.ChunkRequest
}

// EncodeSegment builds the Clinia request from SDK inputs.
func EncodeSegment(texts []string, opts api.EmbeddingOptions) (SegmentParams, error) {
	if len(texts) == 0 {
		return SegmentParams{}, fmt.Errorf("clinia/segment: texts cannot be empty")
	}

	req := cliniaclient.ChunkRequest{
		Texts: texts,
	}

	out := SegmentParams{Request: req}
	return out, nil
}

// DecodeSegment converts the Clinia response into SDK-friendly SegmentingResponse.
func DecodeSegment(resp *cliniaclient.ChunkResponse) (api.SegmentingResponse, error) {
	if resp == nil {
		return api.SegmentingResponse{}, fmt.Errorf("clinia/segment: response is nil")
	}

	segments := make([][]api.Segment, len(resp.Chunks))
	for i, chunkList := range resp.Chunks {
		decoded := make([]api.Segment, len(chunkList))
		for j, chunk := range chunkList {
			decoded[j] = api.Segment{
				ID:         chunk.ID,
				Text:       chunk.Text,
				StartIndex: chunk.StartIndex,
				EndIndex:   chunk.EndIndex,
				TokenCount: chunk.TokenCount,
			}
		}
		segments[i] = decoded
	}

	return api.SegmentingResponse{RequestID: resp.ID, Segments: segments}, nil
}
