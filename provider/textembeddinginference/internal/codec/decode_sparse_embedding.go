package codec

import (
	"fmt"
	"net/http"

	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

// DecodeSparseEmbedding maps the TEI sparse embedding API response to the unified SparseEmbeddingResponse.
// TEI returns sparse embeddings as [][]SparseValue where each SparseValue has index and value
func DecodeSparseEmbedding(resp *tei.CreateSparseEmbeddingResponse) (api.SparseEmbeddingResponse, error) {
	if resp == nil {
		return api.SparseEmbeddingResponse{}, api.NewEmptyResponseBodyError("response from TEI sparse embeddings API is nil")
	}

	// TEI returns [][]SparseValue directly
	// TODO: verify
	sparseArrays := *resp
	sparseEmbs := make([]api.SparseEmbedding, len(sparseArrays))

	for i, sparseValues := range sparseArrays {
		// Convert TEI sparse format to api.SparseEmbedding (map[string]float64)
		sparseMap := make(api.SparseEmbedding)
		for _, sv := range sparseValues {
			// Convert index to string key
			key := fmt.Sprintf("%d", sv.Index)
			sparseMap[key] = sv.Value
		}
		sparseEmbs[i] = sparseMap
	}

	// TEI doesn't return usage information in the basic response
	// Usage would need to be tracked separately if needed
	var usage *api.EmbeddingUsage

	return api.SparseEmbeddingResponse{
		Embeddings: sparseEmbs,
		Usage:      usage,
		RawResponse: &api.EmbeddingRawResponse{
			Headers: http.Header{},
		},
	}, nil
}
