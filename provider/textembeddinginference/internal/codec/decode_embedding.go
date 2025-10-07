package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

// DecodeEmbedding maps the TEI embedding API response to the unified api.EmbeddingResponse.
func DecodeEmbedding(resp *tei.CreateEmbeddingResponse) (api.EmbeddingResponse, error) {
	if resp == nil {
		return api.EmbeddingResponse{}, api.NewEmptyResponseBodyError("response from TEI embeddings API is nil")
	}

	embs := make([]api.Embedding, len(resp.Data))
	for i, d := range resp.Data {
		vec := make([]float64, len(d.Embedding))
		copy(vec, d.Embedding)
		embs[i] = vec
	}

	var usage *api.EmbeddingUsage
	if resp.Usage != nil {
		usage = &api.EmbeddingUsage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}

	return api.EmbeddingResponse{
		Embeddings: embs,
		Usage:      usage,
		RawResponse: &api.EmbeddingRawResponse{
			Headers: http.Header{},
		},
	}, nil
}
