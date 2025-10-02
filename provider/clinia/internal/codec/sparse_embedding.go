package codec

import (
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"go.jetify.com/ai/api"
)

type SparseParams struct {
	Request cliniaclient.SparseEmbedRequest
}

func EncodeSparseEmbedding(texts []string, opts api.EmbeddingOptions) (SparseParams, error) {
	if len(texts) == 0 {
		return SparseParams{}, fmt.Errorf("clinia/sparse: texts cannot be empty")
	}
	out := SparseParams{
		Request: cliniaclient.SparseEmbedRequest{Texts: texts},
	}
	return out, nil
}

func DecodeSparseEmbedding(resp *cliniaclient.SparseEmbedResponse) (api.SparseEmbeddingResponse, error) {
	if resp == nil {
		return api.SparseEmbeddingResponse{}, fmt.Errorf("clinia/sparse: response is nil")
	}
	out := make([]api.SparseEmbedding, len(resp.Embeddings))
	for i, m := range resp.Embeddings {
		conv := make(map[string]float64, len(m))
		for k, v := range m {
			conv[k] = float64(v)
		}
		out[i] = api.SparseEmbedding(conv)
	}
	return api.SparseEmbeddingResponse{Embeddings: out}, nil
}
