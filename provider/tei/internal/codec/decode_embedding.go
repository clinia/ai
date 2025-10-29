package codec

import (
	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/tei/client"
)

// DecodeEmbedding maps the TEI embedding API response to the unified api.EmbeddingResponse.
// TEI returns embeddings as a direct array of arrays: [[0.1, 0.2], [0.3, 0.4]]
func DecodeEmbedding(resp *tei.CreateEmbeddingResponse) (api.DenseEmbeddingResponse, error) {
	if resp == nil {
		return api.DenseEmbeddingResponse{}, api.NewEmptyResponseBodyError("response from TEI embeddings API is nil")
	}

	// TEI returns [][]float64 directly
	embeddingArrays := *resp
	embs := make([]api.Embedding, len(embeddingArrays))
	for i, embedding := range embeddingArrays {
		// Each embedding is already []float64, just copy it
		vec := make([]float64, len(embedding))
		copy(vec, embedding)
		embs[i] = vec
	}

	// TEI doesn't return usage information in the basic response
	// Usage would need to be tracked separately if needed
	var usage *api.EmbeddingUsage

	return api.DenseEmbeddingResponse{
		Embeddings:  embs,
		Usage:       usage,
		RawResponse: nil,
	}, nil
}
