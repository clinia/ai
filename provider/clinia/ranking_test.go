package clinia

import (
	"context"
	"errors"
	"testing"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
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
	ctx := context.Background()

	tests := []struct {
		name     string
		modelID  string
		query    string
		texts    []string
		opts     api.RankingOptions
		ranker   *fakeRanker
		wantErr  bool
		wantResp *api.RankingResponse
		after    func(t *testing.T, ranker *fakeRanker)
	}{
		{
			name:    "successful ranking",
			modelID: "ranker:2",
			query:   "heart",
			texts:   []string{"a", "b"},
			ranker: &fakeRanker{
				response: &cliniaclient.RankResponse{ID: "req", Scores: []float32{0.9, 0.2}},
			},
			wantResp: &api.RankingResponse{RequestID: "req", Scores: []float64{0.9, 0.2}},
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 1, ranker.calls)
				require.Equal(t, "ranker", ranker.lastModelName)
				require.Equal(t, "2", ranker.lastModelVersion)
				require.Equal(t, cliniaclient.RankRequest{
					Query: "heart",
					Texts: []string{"a", "b"},
				}, ranker.lastRequest)
			},
		},
		{
			name:    "propagates provider error",
			modelID: "ranker",
			query:   "q",
			texts:   []string{"a"},
			ranker:  &fakeRanker{err: errors.New("boom")},
			wantErr: true,
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 1, ranker.calls)
			},
		},
		{
			name:    "validates query",
			modelID: "ranker",
			query:   "",
			texts:   []string{"a"},
			ranker:  &fakeRanker{},
			wantErr: true,
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 0, ranker.calls)
			},
		},
		{
			name:    "validates texts",
			modelID: "ranker",
			query:   "q",
			texts:   []string{},
			ranker:  &fakeRanker{},
			wantErr: true,
			after: func(t *testing.T, ranker *fakeRanker) {
				require.Equal(t, 0, ranker.calls)
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			provider, err := NewProvider(ctx, WithRequester(requesterStub{}))
			require.NoError(t, err)
			provider.ranker = tc.ranker

			model, err := provider.NewRankingModel(tc.modelID)
			require.NoError(t, err)
			resp, err := model.Rank(ctx, tc.query, tc.texts, tc.opts)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantResp.RequestID, resp.RequestID)
				require.InDeltaSlice(t, tc.wantResp.Scores, resp.Scores, 1e-6)
			}
			if tc.after != nil {
				tc.after(t, tc.ranker)
			}
		})
	}
}
