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

type chunkerStub struct {
	calls            int
	lastModelName    string
	lastModelVersion string
	lastRequest      cliniaclient.ChunkRequest
	response         *cliniaclient.ChunkResponse
	err              error
	boundRequester   common.Requester
}

func (f *chunkerStub) Chunk(ctx context.Context, modelName, modelVersion string, req cliniaclient.ChunkRequest) (*cliniaclient.ChunkResponse, error) {
	f.calls++
	f.lastModelName = modelName
	f.lastModelVersion = modelVersion
	f.lastRequest = req
	return f.response, f.err
}

func (f *chunkerStub) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func TestChunkingModel(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name              string
		modelName         string
		modelVersion      string
		texts             []string
		opts              api.ChunkingOptions
		chunker           *chunkerStub
		requesterCloseErr error
		wantModelErr      bool
		wantErr           bool
		wantResp          api.ChunkingResponse
		wantModelID       string
		after             func(t *testing.T, chunker *chunkerStub, requester *requesterStub)
	}{
		{
			name:         "successful chunk",
			modelName:    "chunker",
			modelVersion: "2",
			texts:        []string{"hello"},
			chunker: &chunkerStub{
				response: &cliniaclient.ChunkResponse{
					ID: "req",
					Chunks: [][]cliniaclient.Chunk{{
						{ID: "c1", Text: "hello", StartIndex: 0, EndIndex: 5, TokenCount: 5},
					}},
				},
			},
			wantResp: api.ChunkingResponse{
				RequestID: "req",
				Chunks: [][]api.Chunk{{
					{
						ID:         "c1",
						Text:       "hello",
						StartIndex: 0,
						EndIndex:   5,
						TokenCount: 5,
					},
				}},
			},
			wantModelID: "chunker:2",
			after: func(t *testing.T, chunker *chunkerStub, requester *requesterStub) {
				require.Equal(t, 1, chunker.calls)
				require.Equal(t, "chunker", chunker.lastModelName)
				require.Equal(t, "2", chunker.lastModelVersion)
				require.Equal(t, cliniaclient.ChunkRequest{Texts: []string{"hello"}}, chunker.lastRequest)
				require.Equal(t, requester, chunker.boundRequester)
			},
		},
		{
			name:         "propagates provider error",
			modelName:    "chunker",
			modelVersion: "2",
			texts:        []string{"hello"},
			chunker:      &chunkerStub{err: errors.New("boom")},
			wantErr:      true,
			wantModelID:  "chunker:2",
			after: func(t *testing.T, chunker *chunkerStub, requester *requesterStub) {
				require.Equal(t, 1, chunker.calls)
				require.Equal(t, requester, chunker.boundRequester)
			},
		},
		{
			name:         "requires model version",
			modelName:    "chunker",
			modelVersion: "",
			texts:        []string{"hello"},
			chunker:      &chunkerStub{},
			wantModelErr: true,
		},
		{
			name:         "validates texts",
			modelName:    "chunker",
			modelVersion: "2",
			texts:        []string{},
			chunker:      &chunkerStub{},
			wantErr:      true,
			wantModelID:  "chunker:2",
			after: func(t *testing.T, chunker *chunkerStub, requester *requesterStub) {
				require.Equal(t, 0, chunker.calls)
			},
		},
		{
			name:              "close error surfaces",
			modelName:         "chunker",
			modelVersion:      "2",
			texts:             []string{"hello"},
			chunker:           &chunkerStub{response: &cliniaclient.ChunkResponse{}},
			requesterCloseErr: errors.New("close chunk"),
			wantErr:           true,
			wantModelID:       "chunker:2",
			after: func(t *testing.T, chunker *chunkerStub, requester *requesterStub) {
				require.Equal(t, 1, chunker.calls)
				require.Equal(t, requester, chunker.boundRequester)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requester := &requesterStub{closeErr: tt.requesterCloseErr}
			// Inject requester via metadata
			if tt.opts.ProviderMetadata == nil {
				tt.opts.ProviderMetadata = api.NewProviderMetadata(nil)
			}
			tt.opts.ProviderMetadata.Set("clinia", codec.Metadata{Requester: requester})

			provider, err := NewProvider(ctx,
				withChunkerFactory(func(ctx context.Context, opts common.ClientOptions) cliniaclient.Chunker {
					tt.chunker.boundRequester = opts.Requester
					return tt.chunker
				}),
			)
			require.NoError(t, err)

			modelID := tt.modelName
			if tt.modelVersion != "" {
				modelID = tt.modelName + ":" + tt.modelVersion
			}
			model, err := provider.ChunkingModel(modelID)
			if tt.wantModelErr {
				require.Error(t, err)
				require.Equal(t, 0, requester.closeCalls)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantModelID, model.ModelID())

			resp, err := model.Chunk(ctx, tt.texts, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}

			if len(tt.texts) == 0 {
				require.Equal(t, 0, requester.closeCalls)
			} else {
				require.Equal(t, 1, requester.closeCalls)
			}

			if tt.after != nil {
				tt.after(t, tt.chunker, requester)
			}

			if tt.requesterCloseErr != nil && len(tt.texts) > 0 {
				require.ErrorIs(t, err, tt.requesterCloseErr)
			}
		})
	}
}
