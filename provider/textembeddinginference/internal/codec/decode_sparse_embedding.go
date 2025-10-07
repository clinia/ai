package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

// SparseEmbedding represents a sparse embedding vector with indices and values
type SparseEmbedding struct {
	// Indices of non-zero values in the sparse vector
	Indices []int64
	// Values corresponding to the indices
	Values []float64
}

// SparseEmbeddingResponse represents the response from generating sparse embeddings.
type SparseEmbeddingResponse struct {
	// SparseEmbeddings are the generated sparse embeddings. They are in the same order as the input values.
	SparseEmbeddings []SparseEmbedding

	// Usage contains token usage information.
	Usage *api.EmbeddingUsage

	// RawResponse contains optional raw response information for debugging purposes.
	RawResponse *api.EmbeddingRawResponse
}

// DecodeSparseEmbedding maps the TEI sparse embedding API response to the unified SparseEmbeddingResponse.
// TEI returns sparse embeddings as [][]SparseValue where each SparseValue has index and value
func DecodeSparseEmbedding(resp *tei.CreateSparseEmbeddingResponse) (SparseEmbeddingResponse, error) {
	if resp == nil {
		return SparseEmbeddingResponse{}, api.NewEmptyResponseBodyError("response from TEI sparse embeddings API is nil")
	}

	// TEI returns [][]SparseValue directly
	sparseArrays := *resp
	sparseEmbs := make([]SparseEmbedding, len(sparseArrays))

	for i, sparseValues := range sparseArrays {
		indices := make([]int64, len(sparseValues))
		values := make([]float64, len(sparseValues))

		for j, sv := range sparseValues {
			indices[j] = sv.Index
			values[j] = sv.Value
		}

		sparseEmbs[i] = SparseEmbedding{
			Indices: indices,
			Values:  values,
		}
	}

	// TEI doesn't return usage information in the basic response
	// Usage would need to be tracked separately if needed
	var usage *api.EmbeddingUsage

	return SparseEmbeddingResponse{
		SparseEmbeddings: sparseEmbs,
		Usage:            usage,
		RawResponse: &api.EmbeddingRawResponse{
			Headers: http.Header{},
		},
	}, nil
}
