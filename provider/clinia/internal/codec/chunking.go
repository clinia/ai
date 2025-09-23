package codec

import (
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"go.jetify.com/ai/api"
)

// ChunkParams holds the Clinia chunk request.
type ChunkParams struct {
	Request cliniaclient.ChunkRequest
}

// EncodeChunk builds the Clinia chunk request from SDK inputs.
func EncodeChunk(texts []string) (ChunkParams, error) {
	if len(texts) == 0 {
		return ChunkParams{}, fmt.Errorf("clinia/chunk: texts cannot be empty")
	}

	req := cliniaclient.ChunkRequest{
		Texts: texts,
	}

	return ChunkParams{
		Request: req,
	}, nil
}

// DecodeChunk converts the Clinia response into SDK-friendly structures.
func DecodeChunk(resp *cliniaclient.ChunkResponse) (api.ChunkingResponse, error) {
	if resp == nil {
		return api.ChunkingResponse{}, fmt.Errorf("clinia/chunk: response is nil")
	}

	chunks := make([][]api.Chunk, len(resp.Chunks))
	for i, chunkList := range resp.Chunks {
		decoded := make([]api.Chunk, len(chunkList))
		for j, chunk := range chunkList {
			decoded[j] = api.Chunk{
				ID:         chunk.ID,
				Text:       chunk.Text,
				StartIndex: chunk.StartIndex,
				EndIndex:   chunk.EndIndex,
				TokenCount: chunk.TokenCount,
			}
		}
		chunks[i] = decoded
	}

	return api.ChunkingResponse{
		RequestID: resp.ID,
		Chunks:    chunks,
	}, nil
}
