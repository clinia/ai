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

type fakeEmbedder struct {
	lastModelName    string
	lastModelVersion string
	lastRequest      cliniaclient.EmbedRequest
	response         *cliniaclient.EmbedResponse
	err              error
	calls            int
	boundRequester   common.Requester
}

func (f *fakeEmbedder) Embed(ctx context.Context, modelName, modelVersion string, req cliniaclient.EmbedRequest) (*cliniaclient.EmbedResponse, error) {
	f.calls++
	f.lastModelName = modelName
	f.lastModelVersion = modelVersion
	f.lastRequest = req
	return f.response, f.err
}

func (f *fakeEmbedder) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func TestEmbeddingModelDoEmbed(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		values       []string
		embedder     *fakeEmbedder
		baseURL      *string
		wantModelErr bool
		wantErr      bool
		wantResp     *api.DenseEmbeddingResponse
		wantModelID  string
		after        func(t *testing.T, embedder *fakeEmbedder)
	}{
		{
			name:         "successful embedding",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"hello"},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{1, 2}}},
			},
			wantResp: &api.DenseEmbeddingResponse{
				Embeddings: []api.Embedding{
					{1, 2},
				},
			},
			wantModelID: "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
				require.Equal(t, "dense", embedder.lastModelName)
				require.Equal(t, "2", embedder.lastModelVersion)
				require.Equal(t, []string{"hello"}, embedder.lastRequest.Texts)
				require.NotNil(t, embedder.boundRequester)
			},
		},
		{
			name:         "provider returns error",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"hi"},
			embedder:     &fakeEmbedder{err: errors.New("boom")},
			wantErr:      true,
			wantModelID:  "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
				require.NotNil(t, embedder.boundRequester)
			},
		},
		{
			name:         "requires model version",
			modelName:    "dense",
			modelVersion: "",
			values:       []string{"hi"},
			embedder:     &fakeEmbedder{},
			wantModelErr: true,
		},
		{
			name:         "empty values produce validation error",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{}},
			},
			wantErr:     true,
			wantModelID: "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 0, embedder.calls)
			},
		},
		{
			name:         "requester factory error",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"hello"},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{1}}},
			},
			baseURL:     ptr("127.0.0.1"),
			wantErr:     true,
			wantModelID: "dense:2",
			after:       func(t *testing.T, embedder *fakeEmbedder) { require.Equal(t, 0, embedder.calls) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := api.TransportOptions{}
			// Ensure makeRequester is called when inputs are valid
			if tt.baseURL != nil {
				opts.BaseURL = tt.baseURL
			} else if len(tt.values) > 0 {
				host := "127.0.0.1:9000"
				opts.BaseURL = &host
			}

			provider, err := NewProvider(
				withEmbeddingFactory(func(ctx context.Context, opts common.ClientOptions) cliniaclient.Embedder {
					tt.embedder.boundRequester = opts.Requester
					return tt.embedder
				}),
			)
			require.NoError(t, err)

			modelID := tt.modelName
			if tt.modelVersion != "" {
				modelID = tt.modelName + ":" + tt.modelVersion
			}

			model, err := provider.TextEmbeddingModel(modelID)
			if tt.wantModelErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantModelID != "" {
				require.Equal(t, tt.wantModelID, model.ModelID())
			}

			resp, err := model.DoEmbed(ctx, tt.values, opts)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tt.wantResp)
				require.Equal(t, *tt.wantResp, resp)
			}

			if tt.after != nil {
				tt.after(t, tt.embedder)
			}
			if tt.baseURL != nil {
				require.Error(t, err)
			}
		})
	}
}

func ptr(s string) *string { return &s }
