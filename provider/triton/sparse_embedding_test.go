package triton

import (
	"context"
	"errors"
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
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
	boundRequester   common.Requester
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
		name              string
		modelName         string
		modelVersion      string
		texts             []string
		sparse            *sparseStub
		requesterCloseErr error
		wantCtorErr       bool
		wantErr           bool
		want              api.SparseEmbeddingResponse
		after             func(*testing.T, *sparseStub)
	}{
		{
			name:         "successful",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{response: &cliniaclient.SparseEmbedResponse{ID: "req", Embeddings: []map[string]float32{{"x": 0.5}}}},
			want:         api.SparseEmbeddingResponse{Embeddings: []api.SparseEmbedding{{"x": 0.5}}},
			after: func(t *testing.T, s *sparseStub) {
				require.Equal(t, 1, s.calls)
				require.Equal(t, "sparse", s.lastModelName)
				require.Equal(t, "1", s.lastModelVersion)
				require.Equal(t, cliniaclient.SparseEmbedRequest{Texts: []string{"a"}}, s.lastRequest)
				require.NotNil(t, s.boundRequester)
			},
		},
		{
			name:         "provider error",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{err: errors.New("boom")},
			wantErr:      true,
			after: func(t *testing.T, s *sparseStub) {
				require.Equal(t, 1, s.calls)
				require.NotNil(t, s.boundRequester)
			},
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
		// close error surfaces: removed requester injection; ensure standard flow via BaseURL
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure makeRequester is called when texts are provided
			opts := api.TransportOptions{}
			if len(tt.texts) > 0 {
				host := "127.0.0.1:9000"
				opts.BaseURL = &host
			}

			p, err := NewProvider(ctx,
				withSparseFactory(func(ctx context.Context, opts common.ClientOptions) cliniaclient.SparseEmbedder {
					tt.sparse.boundRequester = opts.Requester
					return tt.sparse
				}),
			)
			require.NoError(t, err)

			modelID := tt.modelName
			if tt.modelVersion != "" {
				modelID = tt.modelName + ":" + tt.modelVersion
			}
			m, err := p.SparseEmbeddingModel(modelID)
			if tt.wantCtorErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			resp, err := m.DoEmbed(ctx, tt.texts, opts)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, resp)
			}

			if tt.after != nil {
				tt.after(t, tt.sparse)
			}
		})
	}
}
