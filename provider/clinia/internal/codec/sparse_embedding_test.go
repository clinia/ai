package codec

import (
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeSparseEmbedding(t *testing.T) {
	params, err := EncodeSparseEmbedding("sparse", "1", []string{"a", "b"}, api.SparseEmbeddingOptions{})
	require.NoError(t, err)
	require.Equal(t, "sparse", params.ModelName)
	require.Equal(t, "1", params.ModelVersion)
	require.Equal(t, []string{"a", "b"}, params.Request.Texts)

	_, err = EncodeSparseEmbedding("sparse", "1", nil, api.SparseEmbeddingOptions{})
	require.Error(t, err)
}

func TestDecodeSparseEmbedding(t *testing.T) {
	resp, err := DecodeSparseEmbedding(&cliniaclient.SparseEmbedResponse{
		ID:         "req",
		Embeddings: []map[string]float32{{"token": 1.5}},
	})
	require.NoError(t, err)
	require.Equal(t, "req", resp.RequestID)
	require.Len(t, resp.Embeddings, 1)
	require.Equal(t, 1.5, resp.Embeddings[0]["token"])

	_, err = DecodeSparseEmbedding(nil)
	require.Error(t, err)
}
