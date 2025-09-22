package clinia

import (
	"context"
	"errors"
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
)

type chunkerStub struct {
	calls            int
	lastModelName    string
	lastModelVersion string
	lastRequest      cliniaclient.ChunkRequest
	response         *cliniaclient.ChunkResponse
	err              error
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
	ctx := context.Background()

	tests := []struct {
		name     string
		modelID  string
		texts    []string
		opts     api.ChunkingOptions
		chunker  *chunkerStub
		wantErr  bool
		wantResp api.ChunkingResponse
		after    func(t *testing.T, chunker *chunkerStub)
	}{
		{
			name:    "successful chunk",
			modelID: "chunker:2",
			texts:   []string{"hello"},
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
			after: func(t *testing.T, chunker *chunkerStub) {
				require.Equal(t, 1, chunker.calls)
				require.Equal(t, "chunker", chunker.lastModelName)
				require.Equal(t, "2", chunker.lastModelVersion)
				require.Equal(t, cliniaclient.ChunkRequest{Texts: []string{"hello"}}, chunker.lastRequest)
			},
		},
		{
			name:    "propagates provider error",
			modelID: "chunker",
			texts:   []string{"hello"},
			chunker: &chunkerStub{err: errors.New("boom")},
			wantErr: true,
			after: func(t *testing.T, chunker *chunkerStub) {
				require.Equal(t, 1, chunker.calls)
			},
		},
		{
			name:    "validates texts",
			modelID: "chunker",
			texts:   []string{},
			chunker: &chunkerStub{},
			wantErr: true,
			after: func(t *testing.T, chunker *chunkerStub) {
				require.Equal(t, 0, chunker.calls)
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			provider, err := NewProvider(ctx, WithRequester(requesterStub{}))
			require.NoError(t, err)
			provider.chunker = tc.chunker

			model, err := provider.NewChunkingModel(tc.modelID)
			require.NoError(t, err)

			resp, err := model.Chunk(ctx, tc.texts, tc.opts)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResp, resp)
			}

			if tc.after != nil {
				tc.after(t, tc.chunker)
			}
		})
	}
}
