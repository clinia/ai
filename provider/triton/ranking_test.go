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

type fakeRanker struct {
	lastModelName    string
	lastModelVersion string
	lastRequest      cliniaclient.RankRequest
	response         *cliniaclient.RankResponse
	err              error
	calls            int
	boundRequester   common.Requester
}

func (f *fakeRanker) Rank(ctx context.Context, modelName, modelVersion string, req cliniaclient.RankRequest) (*cliniaclient.RankResponse, error) {
	f.calls++
	f.lastModelName = modelName
	f.lastModelVersion = modelVersion
	f.lastRequest = req
	return f.response, f.err
}

func (f *fakeRanker) Ready(ctx context.Context, modelName, modelVersion string) error { return nil }

func TestRankingModel(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name         string
		modelName    string
		modelVersion string
		query        string
		texts        []string
		opts         api.TransportOptions
		ranker       *fakeRanker
		wantModelErr bool
		wantErr      bool
		wantResp     *api.RankingResponse
		wantModelID  string
		after        func(t *testing.T, ranker *fakeRanker)
	}{
		{
			name:         "successful ranking",
			modelName:    "ranker",
			modelVersion: "2",
			query:        "heart",
			texts:        []string{"a", "b"},
			ranker: &fakeRanker{
				response: &cliniaclient.RankResponse{ID: "req", Scores: []float32{0.9, 0.2}},
			},
			wantResp:    &api.RankingResponse{RequestID: "req", Scores: []float64{0.9, 0.2}},
			wantModelID: "ranker:2",
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 1, ranker.calls)
				require.Equal(t, "ranker", ranker.lastModelName)
				require.Equal(t, "2", ranker.lastModelVersion)
				require.Equal(t, cliniaclient.RankRequest{Query: "heart", Texts: []string{"a", "b"}}, ranker.lastRequest)
				require.NotNil(t, ranker.boundRequester)
			},
		},
		{
			name:         "propagates provider error",
			modelName:    "ranker",
			modelVersion: "2",
			query:        "q",
			texts:        []string{"a"},
			ranker:       &fakeRanker{err: errors.New("boom")},
			wantErr:      true,
			wantModelID:  "ranker:2",
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 1, ranker.calls)
				require.NotNil(t, ranker.boundRequester)
			},
		},
		{
			name:         "requires model version",
			modelName:    "ranker",
			modelVersion: "",
			query:        "q",
			texts:        []string{"a"},
			ranker:       &fakeRanker{},
			wantModelErr: true,
		},
		{
			name:         "validates query",
			modelName:    "ranker",
			modelVersion: "2",
			query:        "",
			texts:        []string{"a"},
			ranker:       &fakeRanker{},
			wantErr:      true,
			wantModelID:  "ranker:2",
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 0, ranker.calls)
			},
		},
		{
			name:         "validates texts",
			modelName:    "ranker",
			modelVersion: "2",
			query:        "q",
			texts:        []string{},
			ranker:       &fakeRanker{},
			wantErr:      true,
			wantModelID:  "ranker:2",
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 0, ranker.calls)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure makeRequester is called when inputs are valid
			if tt.query != "" && len(tt.texts) > 0 {
				host := "127.0.0.1:9000"
				tt.opts.BaseURL = &host
			}

			provider, err := NewProvider(
				withRankerFactory(func(opts common.ClientOptions) cliniaclient.Ranker {
					tt.ranker.boundRequester = opts.Requester
					return tt.ranker
				}),
			)
			require.NoError(t, err)

			modelID := tt.modelName
			if tt.modelVersion != "" {
				modelID = tt.modelName + ":" + tt.modelVersion
			}
			model, err := provider.RankingModel(modelID)
			if tt.wantModelErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantModelID, model.ModelID())

			resp, err := model.DoRank(ctx, tt.query, tt.texts, tt.opts)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tt.wantResp)
				require.Equal(t, tt.wantResp.RequestID, resp.RequestID)
				require.InDeltaSlice(t, tt.wantResp.Scores, resp.Scores, 1e-6)
			}

			if tt.after != nil {
				tt.after(t, tt.ranker)
			}
		})
	}
}
