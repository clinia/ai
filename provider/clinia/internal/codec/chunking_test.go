package codec

import (
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeChunk(t *testing.T) {
	tests := []struct {
		name    string
		texts   []string
		want    cliniaclient.ChunkRequest
		wantErr bool
	}{
		{
			name:  "basic",
			texts: []string{"a", "b"},
			want: cliniaclient.ChunkRequest{
				Texts: []string{"a", "b"},
			},
		},
		{
			name:    "empty texts",
			texts:   []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		parameters := tt
		t.Run(tt.name, func(t *testing.T) {
			params, err := EncodeChunk(parameters.texts)
			if parameters.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, parameters.want, params.Request)
		})
	}
}

func TestDecodeChunk(t *testing.T) {
	tests := []struct {
		name    string
		input   *cliniaclient.ChunkResponse
		want    api.ChunkingResponse
		wantErr bool
	}{
		{
			name: "basic",
			input: &cliniaclient.ChunkResponse{
				ID: "req",
				Chunks: [][]cliniaclient.Chunk{{
					{ID: "c1", Text: "foo", StartIndex: 0, EndIndex: 3, TokenCount: 3},
				}},
			},
			want: api.ChunkingResponse{
				RequestID: "req",
				Chunks: [][]api.Chunk{{
					{
						ID:         "c1",
						Text:       "foo",
						StartIndex: 0,
						EndIndex:   3,
						TokenCount: 3,
					},
				}},
			},
		},
		{
			name:    "nil response",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		parameters := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := DecodeChunk(parameters.input)
			if parameters.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, parameters.want, resp)
		})
	}
}
