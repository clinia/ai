package clinia

import (
	"context"
	"errors"
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

type sparseStub struct {
	calls            int
	lastModelName    string
	lastModelVersion string
	lastRequest      cliniaclient.SparseEmbedRequest
	response         *cliniaclient.SparseEmbedResponse
	err              error
}

func (s *sparseStub) SparseEmbed(ctx context.Context, modelName, modelVersion string, req cliniaclient.SparseEmbedRequest) (*cliniaclient.SparseEmbedResponse, error) {
	s.calls++
	s.lastModelName, s.lastModelVersion, s.lastRequest = modelName, modelVersion, req
	return s.response, s.err
}
func (s *sparseStub) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func TestSparseEmbeddingModel(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		texts        []string
		sparse       *sparseStub
		wantCtorErr  bool
		wantErr      bool
		want         api.SparseEmbeddingResponse
		after        func(*testing.T, *sparseStub)
	}{
		{
			name:         "successful",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{response: &cliniaclient.SparseEmbedResponse{ID: "req", Embeddings: []map[string]float32{{"x": 0.5}}}},
			want:         api.SparseEmbeddingResponse{RequestID: "req", Embeddings: []map[string]float64{{"x": 0.5}}},
			after: func(t *testing.T, s *sparseStub) {
				require.Equal(t, 1, s.calls)
				require.Equal(t, "sparse", s.lastModelName)
				require.Equal(t, "1", s.lastModelVersion)
				require.Equal(t, cliniaclient.SparseEmbedRequest{Texts: []string{"a"}}, s.lastRequest)
			},
		},
		{
			name:         "provider error",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{err: errors.New("boom")},
			wantErr:      true,
			after:        func(t *testing.T, s *sparseStub) { require.Equal(t, 1, s.calls) },
		},
		{
			name:         "constructor needs name",
			modelName:    "",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{},
			wantCtorErr:  true,
		},
		{
			name:         "constructor needs version",
			modelName:    "sparse",
			modelVersion: "",
			texts:        []string{"a"},
			sparse:       &sparseStub{},
			wantCtorErr:  true,
		},
		{
			name:         "validate texts",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{},
			sparse:       &sparseStub{},
			wantErr:      true,
			after:        func(t *testing.T, s *sparseStub) { require.Equal(t, 0, s.calls) },
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewProvider(ctx, WithRequester(requesterStub{}))
			require.NoError(t, err)
			p.sparse = tc.sparse

			modelID := tc.modelName
			if tc.modelVersion != "" {
				modelID = tc.modelName + ":" + tc.modelVersion
			}
			m, err := p.SparseEmbeddingModel(modelID)
			if tc.wantCtorErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			resp, err := m.SparseEmbed(ctx, tc.texts, api.SparseEmbeddingOptions{})
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, resp)
			}

			if tc.after != nil {
				tc.after(t, tc.sparse)
			}
		})
	}
}
