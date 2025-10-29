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

func TestSegmentingModel(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		texts        []string
		opts         api.TransportOptions
		chunker      *chunkerStub
		wantModelErr bool
		wantErr      bool
		wantResp     api.SegmentingResponse
		wantModelID  string
		after        func(t *testing.T, chunker *chunkerStub)
	}{
		{
			name:         "successful segment",
			modelName:    "segmenting",
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
			wantResp: api.SegmentingResponse{
				RequestID: "req",
				Segments: [][]api.Segment{{
					{
						ID:         "c1",
						Text:       "hello",
						StartIndex: 0,
						EndIndex:   5,
						TokenCount: 5,
					},
				}},
			},
			wantModelID: "segmenting:2",
			after: func(t *testing.T, chunker *chunkerStub) {
				require.Equal(t, 1, chunker.calls)
				require.Equal(t, "segmenting", chunker.lastModelName)
				require.Equal(t, "2", chunker.lastModelVersion)
				require.Equal(t, cliniaclient.ChunkRequest{Texts: []string{"hello"}}, chunker.lastRequest)
				require.NotNil(t, chunker.boundRequester)
			},
		},
		{
			name:         "propagates provider error",
			modelName:    "segmenting",
			modelVersion: "2",
			texts:        []string{"hello"},
			chunker:      &chunkerStub{err: errors.New("boom")},
			wantErr:      true,
			wantModelID:  "segmenting:2",
			after: func(t *testing.T, chunker *chunkerStub) {
				require.Equal(t, 1, chunker.calls)
				require.NotNil(t, chunker.boundRequester)
			},
		},
		{
			name:         "requires model version",
			modelName:    "segmenting",
			modelVersion: "",
			texts:        []string{"hello"},
			chunker:      &chunkerStub{},
			wantModelErr: true,
		},
		{
			name:         "validates texts",
			modelName:    "segmenting",
			modelVersion: "2",
			texts:        []string{},
			chunker:      &chunkerStub{},
			wantErr:      true,
			wantModelID:  "segmenting:2",
			after: func(t *testing.T, chunker *chunkerStub) {
				require.Equal(t, 0, chunker.calls)
			},
		},
		// close error surfaces: no longer injectable; ensure normal flow works via BaseURL
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure makeRequester is called when inputs are valid
			if len(tt.texts) > 0 {
				host := "127.0.0.1:9000"
				tt.opts.BaseURL = host
			}

			provider, err := NewProvider(
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
			model, err := provider.SegmentingModel(modelID)
			if tt.wantModelErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantModelID, model.ModelID())

			resp, err := model.DoSegment(ctx, tt.texts, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantResp, resp)
			}

			if tt.after != nil {
				tt.after(t, tt.chunker)
			}
		})
	}
}
