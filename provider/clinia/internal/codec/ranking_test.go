package codec

import (
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

func TestEncodeRank(t *testing.T) {
	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		query        string
		texts        []string
		opts         api.RankingOptions
		want         RankParams
		wantErr      bool
	}{
		{
			name:         "basic",
			modelName:    "ranker",
			modelVersion: "2",
			query:        "heart",
			texts:        []string{"text1", "text2"},
			want: RankParams{
				Request: cliniaclient.RankRequest{
					Query: "heart",
					Texts: []string{"text1", "text2"},
				},
			},
		},
		{
			name:    "empty query",
			query:   "",
			texts:   []string{"body"},
			wantErr: true,
		},
		{
			name:    "empty texts",
			query:   "q",
			texts:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := EncodeRank(tt.query, tt.texts, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, params)
		})
	}
}

func TestDecodeRank(t *testing.T) {
	tests := []struct {
		name    string
		input   *cliniaclient.RankResponse
		want    api.RankingResponse
		wantErr bool
	}{
		{
			name: "basic",
			input: &cliniaclient.RankResponse{
				ID:     "abc",
				Scores: []float32{0.9, 0.1},
			},
			want: api.RankingResponse{
				Scores:    []float64{0.9, 0.1},
				RequestID: "abc",
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
			resp, err := DecodeRank(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want.RequestID, resp.RequestID)
			require.InDeltaSlice(t, tt.want.Scores, resp.Scores, 1e-6)
		})
	}
}
