package codec

import (
	"fmt"
	"strings"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"go.jetify.com/ai/api"
)

// EmbeddingParams represents a fully prepared request for the Clinia embedder.
type EmbeddingParams struct {
	ModelName    string
	ModelVersion string
	Request      cliniaclient.EmbedRequest
}

// EncodeEmbedding converts the SDK request into the Clinia embedder request payload.
func EncodeEmbedding(modelName, modelVersion string, values []string, opts api.EmbeddingOptions) (EmbeddingParams, error) {
	if len(values) == 0 {
		return EmbeddingParams{}, fmt.Errorf("clinia/embed: values cannot be empty")
	}
	if strings.TrimSpace(modelName) == "" {
		return EmbeddingParams{}, fmt.Errorf("clinia/embed: model name is required")
	}

	params := EmbeddingParams{
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Request: cliniaclient.EmbedRequest{
			Texts: values,
		},
	}

	return params, nil
}

// DecodeEmbedding transforms the Clinia embedder response into the SDK shape.
func DecodeEmbedding(resp *cliniaclient.EmbedResponse) (api.EmbeddingResponse, error) {
	if resp == nil {
		return api.EmbeddingResponse{}, fmt.Errorf("clinia/embed: response is nil")
	}

	embeddings := make([]api.Embedding, len(resp.Embeddings))
	for i, embedding := range resp.Embeddings {
		converted := make(api.Embedding, len(embedding))
		for j, value := range embedding {
			converted[j] = float64(value)
		}
		embeddings[i] = converted
	}

	return api.EmbeddingResponse{Embeddings: embeddings}, nil
}
