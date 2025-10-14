package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
)

// DecodeEmbedding maps the Jina embedding API response to the unified api.EmbeddingResponse.
func DecodeEmbedding(resp *jina.CreateEmbeddingResponse) (api.DenseEmbeddingResponse, error) {
	if resp == nil {
		return api.DenseEmbeddingResponse{}, api.NewEmptyResponseBodyError("response from Jina embeddings API is nil")
	}

	embs := make([]api.Embedding, len(resp.Data))
	for i, d := range resp.Data {
		vec := make([]float64, len(d.Embedding))
		copy(vec, d.Embedding)
		embs[i] = vec
	}

	usage := &api.EmbeddingUsage{
		PromptTokens: resp.Usage.PromptTokens,
		TotalTokens:  resp.Usage.TotalTokens,
	}

	return api.DenseEmbeddingResponse{
		Embeddings: embs,
		Usage:      usage,
		RawResponse: &api.EmbeddingRawResponse{
			Headers: http.Header{},
		},
	}, nil
}
