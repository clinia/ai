package codec

import (
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"go.jetify.com/ai/api"
)

type SparseParams struct {
	ModelName    string
	ModelVersion string
	Request      cliniaclient.SparseEmbedRequest
}

func EncodeSparseEmbedding(modelName, modelVersion string, texts []string, _ api.SparseEmbeddingOptions) (SparseParams, error) {
	if len(texts) == 0 {
		return SparseParams{}, fmt.Errorf("clinia/sparse: texts cannot be empty")
	}
	return SparseParams{
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Request:      cliniaclient.SparseEmbedRequest{Texts: texts},
	}, nil
}

func DecodeSparseEmbedding(resp *cliniaclient.SparseEmbedResponse) (api.SparseEmbeddingResponse, error) {
	if resp == nil {
		return api.SparseEmbeddingResponse{}, fmt.Errorf("clinia/sparse: response is nil")
	}
	out := make([]map[string]float64, len(resp.Embeddings))
	for i, m := range resp.Embeddings {
		conv := make(map[string]float64, len(m))
		for k, v := range m {
			conv[k] = float64(v)
		}
		out[i] = conv
	}
	return api.SparseEmbeddingResponse{RequestID: resp.ID, Embeddings: out}, nil
}
