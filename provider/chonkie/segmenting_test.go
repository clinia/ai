package chonkie

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	chonkieClient "go.jetify.com/ai/provider/chonkie/client"
	"go.jetify.com/ai/provider/internal/requesterx"
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
			name:  "batched request body",
			texts: []string{"Hello", "World"},
			exchanges: []httpmock.Exchange{
				{
					Request:  httpmock.Request{Method: http.MethodPost, Path: "/", Body: `{"content":["Hello","World"],"return_chunks":true}`},
					Response: httpmock.Response{StatusCode: http.StatusOK, Body: `[{"num_tokens":2,"tokenizer":"t","num_chunks":1,"chunk_positions":[[0,5]],"chunks":["Hello"]},{"num_tokens":1,"tokenizer":"t","num_chunks":1,"chunk_positions":[[0,5]],"chunks":["World"]}]`},
				},
			},
			expectedResp: api.SegmentingResponse{Segments: [][]api.Segment{{{Text: "Hello"}}, {{Text: "World"}}}},
		},
		{
			name:      "empty input returns error",
			texts:     []string{},
			exchanges: []httpmock.Exchange{},
			wantErr:   true,
		},
		{
			name:  "provider returns HTTP error",
			texts: []string{"Oops"},
			exchanges: []httpmock.Exchange{
				{
					Request:  httpmock.Request{Method: http.MethodPost, Path: "/", Body: `{"content":["Oops"],"return_chunks":true}`},
					Response: httpmock.Response{StatusCode: http.StatusInternalServerError, Body: `{"error":"boom"}`},
				},
			},
			wantErr: true,
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

			client := chonkieClient.NewClient(
				requesterx.WithBaseURL(server.BaseURL()),
				requesterx.WithAPIKey("test-key"),
			)

			provider := NewProvider(WithClient(client))
			model, err := provider.SegmentingModel("segmenting:1")
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
