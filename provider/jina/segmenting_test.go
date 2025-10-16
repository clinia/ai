package jina

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	jinaClient "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/client/option"
	"go.jetify.com/pkg/httpmock"
)

func TestDoSegment(t *testing.T) {
	tests := []struct {
		name         string
		texts        []string
		options      api.TransportOptions
		exchanges    []httpmock.Exchange
		wantErr      bool
		expectedResp api.SegmentingResponse
		skip         bool
	}{
		{
			name:  "default per-text requests",
			texts: []string{"Hello", "World"},
			exchanges: []httpmock.Exchange{
				{
					Request:  httpmock.Request{Method: http.MethodPost, Path: "/segment", Body: `{"content":"Hello","return_chunks":true}`},
					Response: httpmock.Response{StatusCode: http.StatusOK, Body: `{"num_tokens":3,"tokenizer":"t","num_chunks":2,"chunk_positions":[[0,5],[6,11]],"chunks":["Hello","World!"]}`},
				},
				{
					Request:  httpmock.Request{Method: http.MethodPost, Path: "/segment", Body: `{"content":"World","return_chunks":true}`},
					Response: httpmock.Response{StatusCode: http.StatusOK, Body: `{"num_tokens":1,"tokenizer":"t","num_chunks":1,"chunk_positions":[[0,5]],"chunks":["World"]}`},
				},
			},
			expectedResp: api.SegmentingResponse{Segments: [][]api.Segment{{{Text: "Hello"}, {Text: "World!"}}, {{Text: "World"}}}},
		},
	}

	runDoSegmentTests(t, tests)
}

func runDoSegmentTests(t *testing.T, tests []struct {
	name         string
	texts        []string
	options      api.TransportOptions
	exchanges    []httpmock.Exchange
	wantErr      bool
	expectedResp api.SegmentingResponse
	skip         bool
},
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("Skipping test: %s", tt.name)
			}

			server := httpmock.NewServer(t, tt.exchanges)
			defer server.Close()

			client := jinaClient.NewClient(
				option.WithBaseURL(server.BaseURL()),
				option.WithAPIKey("test-key"),
			)

			provider := NewProvider(WithClient(client))
			model, err := provider.Segmenter("segmenter:1")
			require.NoError(t, err)

			resp, err := model.DoSegment(t.Context(), tt.texts, tt.options)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Compare only the texts to avoid being brittle with positions
			gotTexts := make([][]string, len(resp.Segments))
			for i := range resp.Segments {
				row := make([]string, len(resp.Segments[i]))
				for j := range resp.Segments[i] {
					row[j] = resp.Segments[i][j].Text
				}
				gotTexts[i] = row
			}
			wantTexts := make([][]string, len(tt.expectedResp.Segments))
			for i := range tt.expectedResp.Segments {
				row := make([]string, len(tt.expectedResp.Segments[i]))
				for j := range tt.expectedResp.Segments[i] {
					row[j] = tt.expectedResp.Segments[i][j].Text
				}
				wantTexts[i] = row
			}
			require.Equal(t, wantTexts, gotTexts)
		})
	}
}
