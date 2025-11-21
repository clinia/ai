package codec

import (
	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/tei/client"
)

// DecodeSparseEmbedding maps the TEI sparse embedding API response to the unified SparseEmbeddingResponse.
// TEI returns sparse embeddings as [][]SparseValue where each SparseValue has index and value
func DecodeSparseEmbedding(resp *tei.CreateSparseEmbeddingResponse) (api.SparseEmbeddingResponse, error) {
	if resp == nil {
		return api.SparseEmbeddingResponse{}, api.NewEmptyResponseBodyError("response from TEI sparse embeddings API is nil")
	}

	// TEI returns [][]SparseValue directly
	sparseArrays := *resp
	sparseEmbs := make([]api.SparseEmbedding, len(sparseArrays))

	for i, sparseValues := range sparseArrays {
		// Convert TEI sparse format to api.SparseEmbedding (map[string]float64)
		sparseEmbs[i] = api.SparseEmbedding(sparseValues)
	}

	// TEI doesn't return usage information in the basic response
	// Usage would need to be tracked separately if needed
	var usage *api.EmbeddingUsage

	return api.SparseEmbeddingResponse{
		Embeddings:  sparseEmbs,
		Usage:       usage,
		RawResponse: nil,
	}, nil
}
