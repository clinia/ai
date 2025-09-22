package clinia

import (
	"context"
	"errors"
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
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
}

func (f *fakeEmbedder) Embed(ctx context.Context, modelName, modelVersion string, req cliniaclient.EmbedRequest) (*cliniaclient.EmbedResponse, error) {
	f.calls++
	f.lastModelName = modelName
	f.lastModelVersion = modelVersion
	f.lastRequest = req
	return f.response, f.err
}

func (f *fakeEmbedder) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func TestNewProvider(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		values       []string
		embedder     *fakeEmbedder
		wantModelErr bool
		wantErr      bool
		wantResp     *api.EmbeddingResponse
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
			wantModelErr: false,
			wantResp: &api.EmbeddingResponse{
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
			},
		},
		{
			name:         "provider returns error",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"hi"},
			embedder: &fakeEmbedder{
				err: errors.New("boom"),
			},
			wantModelErr: false,
			wantErr:      true,
			wantModelID:  "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
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
			wantModelErr: false,
			wantErr:      true,
			wantModelID:  "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 0, embedder.calls)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(ctx, tt.opts...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, provider)
			if tt.assert != nil {
				tt.assert(t, provider)
			}
		})
	}
}

func TestEmbeddingModelDoEmbed(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		values       []string
		embedder     *fakeEmbedder
		wantModelErr bool
		wantErr      bool
		wantResp     *api.EmbeddingResponse
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
			wantModelErr: false,
			wantResp: &api.EmbeddingResponse{
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
			},
		},
		{
			name:         "provider returns error",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"hi"},
			embedder:     &fakeEmbedder{err: errors.New("boom")},
			wantModelErr: false,
			wantErr:      true,
			wantModelID:  "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
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
			wantModelErr: false,
			wantErr:      true,
			wantModelID:  "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 0, embedder.calls)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(ctx, WithRequester(requesterStub{}))
			require.NoError(t, err)
			provider.embedder = tt.embedder

			model, err := provider.NewEmbeddingModel(tt.modelName, tt.modelVersion)
			if tt.wantModelErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantModelID != "" {
				require.Equal(t, tt.wantModelID, model.ModelID())
			}

			resp, err := model.DoEmbed(ctx, tt.values, api.EmbeddingOptions{})
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
		})
	}
}
