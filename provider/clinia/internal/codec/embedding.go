package codec

import (
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
	"go.jetify.com/ai/api"
)

// EmbeddingParams represents a fully prepared request for the Clinia embedder.
type EmbeddingParams struct {
	Request   cliniaclient.EmbedRequest
	Requester common.Requester
}

// EncodeEmbedding converts the SDK request into the Clinia embedder request payload.
func EncodeEmbedding(values []string, opts api.EmbeddingOptions) (EmbeddingParams, error) {
	if len(values) == 0 {
		return EmbeddingParams{}, fmt.Errorf("clinia/embed: values cannot be empty")
	}

	params := EmbeddingParams{
		Request: cliniaclient.EmbedRequest{
			Texts: values,
		},
	}

	if meta := GetMetadata(opts); meta != nil && meta.Requester != nil {
		params.Requester = meta.Requester
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
