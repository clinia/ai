package codec

import (
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeSparseEmbedding(t *testing.T) {
	params, err := EncodeSparseEmbedding([]string{"a", "b"}, api.EmbeddingOptions{})
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b"}, params.Request.Texts)

	_, err = EncodeSparseEmbedding(nil, api.EmbeddingOptions{})
	require.Error(t, err)
}

func TestDecodeSparseEmbedding(t *testing.T) {
	resp, err := DecodeSparseEmbedding(&cliniaclient.SparseEmbedResponse{
		ID:         "req",
		Embeddings: []map[string]float32{{"token": 1.5}},
	})
	require.NoError(t, err)
	require.Len(t, resp.Embeddings, 1)
	require.Equal(t, 1.5, resp.Embeddings[0]["token"])

	_, err = DecodeSparseEmbedding(nil)
	require.Error(t, err)
}
