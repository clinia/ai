package clinia

import (
	"context"
	"errors"
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
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
		after             func(*testing.T, *sparseStub, *requesterStub)
	}{
		{
			name:         "successful",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{response: &cliniaclient.SparseEmbedResponse{ID: "req", Embeddings: []map[string]float32{{"x": 0.5}}}},
			want:         api.SparseEmbeddingResponse{RequestID: "req", Embeddings: []map[string]float64{{"x": 0.5}}},
			after: func(t *testing.T, s *sparseStub, r *requesterStub) {
				require.Equal(t, 1, s.calls)
				require.Equal(t, "sparse", s.lastModelName)
				require.Equal(t, "1", s.lastModelVersion)
				require.Equal(t, cliniaclient.SparseEmbedRequest{Texts: []string{"a"}}, s.lastRequest)
				require.Equal(t, r, s.boundRequester)
			},
		},
		{
			name:         "provider error",
			modelName:    "sparse",
			modelVersion: "1",
			texts:        []string{"a"},
			sparse:       &sparseStub{err: errors.New("boom")},
			wantErr:      true,
			after: func(t *testing.T, s *sparseStub, r *requesterStub) {
				require.Equal(t, 1, s.calls)
				require.Equal(t, r, s.boundRequester)
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
			after: func(t *testing.T, s *sparseStub, r *requesterStub) {
				require.Equal(t, 0, s.calls)
			},
		},
		{
			name:              "close error surfaces",
			modelName:         "sparse",
			modelVersion:      "1",
			texts:             []string{"a"},
			sparse:            &sparseStub{response: &cliniaclient.SparseEmbedResponse{}},
			requesterCloseErr: errors.New("close sparse"),
			wantErr:           true,
			after: func(t *testing.T, s *sparseStub, r *requesterStub) {
				require.Equal(t, 1, s.calls)
				require.Equal(t, r, s.boundRequester)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requester := &requesterStub{closeErr: tt.requesterCloseErr}
			// Inject requester via metadata
			opts := api.SparseEmbeddingOptions{}
			opts.ProviderMetadata = api.NewProviderMetadata(nil)
			opts.ProviderMetadata.Set("clinia", codec.Metadata{Requester: requester})

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
				require.Equal(t, 0, requester.closeCalls)
				return
			}
			require.NoError(t, err)

			resp, err := m.SparseEmbed(ctx, tt.texts, opts)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, resp)
			}

			if len(tt.texts) == 0 {
				require.Equal(t, 0, requester.closeCalls)
			} else {
				require.Equal(t, 1, requester.closeCalls)
			}

			if tt.after != nil {
				tt.after(t, tt.sparse, requester)
			}

			if tt.requesterCloseErr != nil && len(tt.texts) > 0 {
				require.ErrorIs(t, err, tt.requesterCloseErr)
			}
		})
	}
}
