package codec

import (
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeEmbedding(t *testing.T) {
	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		values       []string
		want         EmbeddingParams
		wantErr      bool
	}{
		{
			name:         "basic",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"foo", "bar"},
			want: EmbeddingParams{
				Request: cliniaclient.EmbedRequest{Texts: []string{"foo", "bar"}},
			},
		},
		{
			name:         "default version",
			modelName:    "dense",
			modelVersion: "",
			values:       []string{"foo"},
			want: EmbeddingParams{
				Request: cliniaclient.EmbedRequest{Texts: []string{"foo"}},
			},
		},
		{
			name:    "empty values",
			values:  []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := EncodeEmbedding(tt.values, api.EmbeddingOptions{})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, params)
		})
	}
}

func TestDecodeEmbedding(t *testing.T) {
	tests := []struct {
		name    string
		input   *cliniaclient.EmbedResponse
		want    api.EmbeddingResponse
		wantErr bool
	}{
		{
			name:  "basic",
			input: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{1.5, 2.5}}},
			want: api.EmbeddingResponse{
				Embeddings: []api.Embedding{
					{1.5, 2.5},
				},
			},
		},
		{
			name:    "nil response",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := DecodeEmbedding(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
