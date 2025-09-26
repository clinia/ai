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
		name              string
		modelName         string
		modelVersion      string
		values            []string
		embedder          *fakeEmbedder
		requesterCloseErr error
		baseURL           *string
		wantModelErr      bool
		wantErr           bool
		wantResp          *api.EmbeddingResponse
		wantModelID       string
		after             func(t *testing.T, embedder *fakeEmbedder, requester *requesterStub)
	}{
		{
			name:         "successful embedding",
			modelName:    "dense",
			modelVersion: "2",
			values:       []string{"hello"},
			embedder: &fakeEmbedder{
				response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{1, 2}}},
			},
			wantResp: &api.EmbeddingResponse{
				Embeddings: []api.Embedding{
					{1, 2},
				},
			},
			wantModelID: "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder, requester *requesterStub) {
				require.Equal(t, 1, embedder.calls)
				require.Equal(t, "dense", embedder.lastModelName)
				require.Equal(t, "2", embedder.lastModelVersion)
				require.Equal(t, []string{"hello"}, embedder.lastRequest.Texts)
				require.Equal(t, requester, embedder.boundRequester)
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
			after: func(t *testing.T, embedder *fakeEmbedder, requester *requesterStub) {
				require.Equal(t, 1, embedder.calls)
				require.Equal(t, requester, embedder.boundRequester)
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
			after: func(t *testing.T, embedder *fakeEmbedder, requester *requesterStub) {
				require.Equal(t, 0, embedder.calls)
			},
		},
		{
			name:              "close error surfaces",
			modelName:         "dense",
			modelVersion:      "2",
			values:            []string{"hello"},
			embedder:          &fakeEmbedder{response: &cliniaclient.EmbedResponse{Embeddings: [][]float32{{3}}}},
			requesterCloseErr: errors.New("close boom"),
			wantErr:           true,
			wantModelID:       "dense:2",
			after: func(t *testing.T, embedder *fakeEmbedder, requester *requesterStub) {
				require.Equal(t, 1, embedder.calls)
				require.Equal(t, requester, embedder.boundRequester)
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
			after: func(t *testing.T, embedder *fakeEmbedder, requester *requesterStub) {
				require.Equal(t, 0, embedder.calls)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requester := &requesterStub{closeErr: tt.requesterCloseErr}
			// Build options: either inject requester or force BaseURL path for error
			opts := api.EmbeddingOptions{}
			if tt.baseURL != nil {
				opts.BaseURL = tt.baseURL
			} else {
				if opts.ProviderMetadata == nil {
					opts.ProviderMetadata = api.NewProviderMetadata(nil)
				}
				opts.ProviderMetadata.Set("clinia", codec.Metadata{Requester: requester})
			}

			provider, err := NewProvider(ctx,
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
				require.Equal(t, 0, requester.closeCalls)
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

			if tt.baseURL != nil {
				// invalid baseURL path should not close
				require.Equal(t, 0, requester.closeCalls)
			} else if len(tt.values) == 0 {
				require.Equal(t, 0, requester.closeCalls)
			} else {
				require.Equal(t, 1, requester.closeCalls)
			}

			if tt.after != nil {
				tt.after(t, tt.embedder, requester)
			}
			if tt.requesterCloseErr != nil && tt.baseURL == nil && len(tt.values) > 0 {
				require.ErrorIs(t, err, tt.requesterCloseErr)
			}
			if tt.baseURL != nil {
				require.Error(t, err)
			}
		})
	}
}

func ptr(s string) *string { return &s }
