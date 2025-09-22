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
		name    string
		opts    []Option
		wantErr bool
		assert  func(t *testing.T, provider *Provider)
	}{
		{
			name:    "requires requester",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "creates embedder with requester",
			opts: []Option{
				WithRequester(requesterStub{}),
			},
			assert: func(t *testing.T, provider *Provider) {
				require.Equal(t, "clinia", provider.Name())
				require.NotNil(t, provider.Embedder())
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
		name     string
		modelID  string
		values   []string
		embedder *fakeEmbedder
		wantErr  bool
		wantResp *api.EmbeddingResponse
		after    func(t *testing.T, embedder *fakeEmbedder)
	}{
		{
			name:    "successful embedding",
			modelID: "dense:2",
			values:  []string{"hello"},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{1, 2}}},
			},
			wantResp: &api.EmbeddingResponse{
				Embeddings: []api.Embedding{
					{1, 2},
				},
			},
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
				require.Equal(t, "dense", embedder.lastModelName)
				require.Equal(t, "2", embedder.lastModelVersion)
				require.Equal(t, []string{"hello"}, embedder.lastRequest.Texts)
			},
		},
		{
			name:    "provider returns error",
			modelID: "dense:2",
			values:  []string{"hi"},
			embedder: &fakeEmbedder{
				err: errors.New("boom"),
			},
			wantErr: true,
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
			},
		},
		{
			name:    "defaults model version when omitted",
			modelID: "dense",
			values:  []string{"hi"},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{3}}},
			},
			wantResp: &api.EmbeddingResponse{
				Embeddings: []api.Embedding{{3}},
			},
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 1, embedder.calls)
				require.Equal(t, "dense", embedder.lastModelName)
				require.Equal(t, "1", embedder.lastModelVersion)
			},
		},
		{
			name:    "empty values produce validation error",
			modelID: "dense:2",
			values:  []string{},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{}},
			},
			wantErr: true,
			after: func(t *testing.T, embedder *fakeEmbedder) {
				require.Equal(t, 0, embedder.calls)
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			provider, err := NewProvider(ctx, WithRequester(requesterStub{}))
			require.NoError(t, err)
			provider.embedder = tc.embedder

			model, err := provider.NewEmbeddingModel(tc.modelID)
			require.NoError(t, err)
			resp, err := model.DoEmbed(ctx, tc.values, api.EmbeddingOptions{})
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tc.wantResp)
				require.Equal(t, *tc.wantResp, resp)
			}

			if tc.after != nil {
				tc.after(t, tc.embedder)
			}
		})
	}
}
