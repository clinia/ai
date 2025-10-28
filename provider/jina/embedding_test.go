package jina

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/client/option"
	"go.jetify.com/pkg/httpmock"
)

func TestDoEmbed(t *testing.T) {
	standardInput := []string{"Hello", "World"}

	// Standard Jina response body used across tests
	standardResponseBody := `{
        "object": "list",
        "data": [
            {
                "object": "embedding",
                "embedding": [0.0023064255, -0.009327292, 0.015797527],
                "index": 0
            },
            {
                "object": "embedding",
                "embedding": [0.0072664247, -0.008545238, 0.017125098],
                "index": 1
            }
        ],
        "model": "text-embedding-ada-002",
        "usage": {
            "prompt_tokens": 2,
            "total_tokens": 2
        }
    }`

	standardExchange := []httpmock.Exchange{
		{
			Request: httpmock.Request{
				Method: http.MethodPost,
				Path:   "/embeddings",
				Body: `{
                    "input": ["Hello", "World"],
                    "model": "text-embedding-ada-002"
                }`,
			},
			Response: httpmock.Response{
				StatusCode: http.StatusOK,
				Body:       standardResponseBody,
			},
		},
	}

	tests := []struct {
		name         string
		modelID      string
		input        []string
		options      api.TransportOptions
		exchanges    []httpmock.Exchange
		wantErr      bool
		expectedResp api.DenseEmbeddingResponse
		skip         bool
	}{
		{
			name:      "successful embedding",
			modelID:   "text-embedding-ada-002",
			input:     standardInput,
			exchanges: standardExchange,
			expectedResp: api.DenseEmbeddingResponse{
				Embeddings: []api.Embedding{
					{0.0023064255, -0.009327292, 0.015797527},
					{0.0072664247, -0.008545238, 0.017125098},
				},
				Usage: &api.EmbeddingUsage{
					PromptTokens: 2,
					TotalTokens:  2,
				},
				RawResponse: &api.EmbeddingRawResponse{
					Headers: http.Header{},
				},
			},
		},
		{
			name:    "with custom headers",
			modelID: "text-embedding-ada-002",
			input:   standardInput,
			options: api.TransportOptions{
				Headers: http.Header{
					"Custom-Header": []string{"test-value"},
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/embeddings",
						Headers: map[string]string{
							"Custom-Header": "test-value",
						},
						Body: `{
                            "input": ["Hello", "World"],
                            "model": "text-embedding-ada-002"
                        }`,
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body:       standardResponseBody,
					},
				},
			},
			expectedResp: api.DenseEmbeddingResponse{
				Embeddings: []api.Embedding{
					{0.0023064255, -0.009327292, 0.015797527},
					{0.0072664247, -0.008545238, 0.017125098},
				},
				Usage: &api.EmbeddingUsage{
					PromptTokens: 2,
					TotalTokens:  2,
				},
				RawResponse: &api.EmbeddingRawResponse{
					Headers: http.Header{},
				},
			},
		},
	}

	runDoEmbedTests(t, tests)
}

func runDoEmbedTests(t *testing.T, tests []struct {
	name         string
	modelID      string
	input        []string
	options      api.TransportOptions
	exchanges    []httpmock.Exchange
	wantErr      bool
	expectedResp api.DenseEmbeddingResponse
	skip         bool
},
) {
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.skip {
				t.Skipf("Skipping test: %s", testCase.name)
			}

			server := httpmock.NewServer(t, testCase.exchanges)
			defer server.Close()

			// Set up client options for the Jina client
			clientOptions := []option.RequestOption{
				option.WithBaseURL(server.BaseURL()),
				option.WithAPIKey("test-key"),
			}

			// Create client with options
			client := jina.NewClient(clientOptions...)

			// Use custom model ID
			modelID := testCase.modelID

			// instantiate the provider with the mocked client
			provider := NewProvider(WithClient(client))

			// Create model with the provider
			model, _ := provider.TextEmbeddingModel(modelID)

			// Build embedding options
			resp, err := model.DoEmbed(t.Context(), testCase.input, testCase.options)

			if testCase.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// Compare response fields
			require.Equal(t, testCase.expectedResp, resp)
		})
	}
}
