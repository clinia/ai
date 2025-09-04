package codec

import (
	"net/http"
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
)

func ptrString(s string) *string {
	return &s
}

func TestEncodeEmbedding(t *testing.T) {
	t.Run("text embedding", func(t *testing.T) {
		tests := []struct {
			name            string
			modelID         string
			values          []string
			headers         http.Header
			wantReqOptsLen  int
			wantWarningsLen int
			expectedParams  jina.TextEmbeddingNewParams
		}{
			{
				name:            "no headers, two values",
				modelID:         "text-embedding-3-small",
				values:          []string{"hello", "world"},
				headers:         nil,
				wantReqOptsLen:  0,
				wantWarningsLen: 0,
				expectedParams: jina.TextEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-small"),
					Input: []string{"hello", "world"},
				},
			},
			{
				name:    "with single and multi-value headers",
				modelID: "text-embedding-3-small",
				values:  []string{"a", "b", "c"},
				headers: func() http.Header {
					h := http.Header{}
					h.Add("X-One", "1")
					h.Add("X-Multi", "A")
					h.Add("X-Multi", "B")
					return h
				}(),
				// 1 option for X-One + 2 options for X-Multi
				wantReqOptsLen:  3,
				wantWarningsLen: 0,
				expectedParams: jina.TextEmbeddingNewParams{
					Model: openai.EmbeddingModel("text-embedding-3-small"),
					Input: []string{"a", "b", "c"},
				},
			},
			{
				name:            "empty input slice",
				modelID:         "text-embedding-3-large",
				values:          []string{},
				headers:         nil,
				wantReqOptsLen:  0,
				wantWarningsLen: 0,
				expectedParams: jina.TextEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-large"),
					Input: []string{},
				},
			},
			{
				name:            "different model id",
				modelID:         "text-embedding-3-small",
				values:          []string{"only one"},
				headers:         http.Header{},
				wantReqOptsLen:  0,
				wantWarningsLen: 0,
				expectedParams: jina.TextEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-small"),
					Input: []string{"only one"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := api.EmbeddingOptions{Headers: tt.headers}

				params, reqOpts, warnings, err := EncodeEmbedding(tt.modelID, tt.values, opts)
				require.NoError(t, err)

				// Request options (opaque type): assert count derived from headers
				assert.Len(t, reqOpts, tt.wantReqOptsLen)

				// Warnings (currently none expected)
				assert.Len(t, warnings, tt.wantWarningsLen)

				// Params: model id
				assert.Equal(t, openai.EmbeddingModel(tt.modelID), params.Model)

				// Params: input union mirrors provided values
				assert.Equal(t, tt.values, params.Input)
			})
		}
	})
	t.Run("multimodal embedding", func(t *testing.T) {
		tests := []struct {
			name            string
			modelID         string
			values          []jina.MultimodalEmbeddingInput
			headers         http.Header
			wantReqOptsLen  int
			wantWarningsLen int
			expectedParams  jina.MultimodalEmbeddingNewParams
		}{
			{
				name:    "no headers, two values",
				modelID: "text-embedding-3-small",
				values: []jina.MultimodalEmbeddingInput{
					{Text: ptrString("hello")}, {Text: ptrString("world")}, {Image: ptrString("")},
				},
				headers:         nil,
				wantReqOptsLen:  0,
				wantWarningsLen: 0,
				expectedParams: jina.MultimodalEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-small"),
					Input: []jina.MultimodalEmbeddingInput{
						{Text: ptrString("hello")},
						{Text: ptrString("world")},
						{Image: ptrString("")},
					},
				},
			},
			{
				name:    "with single and multi-value headers",
				modelID: "text-embedding-3-small",
				values: []jina.MultimodalEmbeddingInput{
					{Text: ptrString("a")},
					{Image: ptrString("alsdkjfa")},
					{Text: ptrString("b")},
				},
				headers: func() http.Header {
					h := http.Header{}
					h.Add("X-One", "1")
					h.Add("X-Multi", "A")
					h.Add("X-Multi", "B")
					return h
				}(),
				// 1 option for X-One + 2 options for X-Multi
				wantReqOptsLen:  3,
				wantWarningsLen: 0,
				expectedParams: jina.MultimodalEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-small"),
					Input: []jina.MultimodalEmbeddingInput{
						{Text: ptrString("a")},
						{Image: ptrString("alsdkjfa")},
						{Text: ptrString("b")},
					},
				},
			},
			{
				name:            "empty input slice",
				modelID:         "text-embedding-3-large",
				values:          []jina.MultimodalEmbeddingInput{},
				headers:         nil,
				wantReqOptsLen:  0,
				wantWarningsLen: 0,
				expectedParams: jina.MultimodalEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-large"),
					Input: []jina.MultimodalEmbeddingInput{},
				},
			},
			{
				name:            "different model id",
				modelID:         "text-embedding-3-small",
				values:          []jina.MultimodalEmbeddingInput{{Text: ptrString("only one")}},
				headers:         http.Header{},
				wantReqOptsLen:  0,
				wantWarningsLen: 0,
				expectedParams: jina.MultimodalEmbeddingNewParams{
					Model: jina.EmbeddingModel("text-embedding-3-small"),
					Input: []jina.MultimodalEmbeddingInput{{Text: ptrString("only one")}},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				opts := api.EmbeddingOptions{Headers: tt.headers}

				params, reqOpts, warnings, err := EncodeMultimodalEmbedding(tt.modelID, tt.values, opts)
				require.NoError(t, err)

				// Request options (opaque type): assert count derived from headers
				assert.Len(t, reqOpts, tt.wantReqOptsLen)

				// Warnings (currently none expected)
				assert.Len(t, warnings, tt.wantWarningsLen)

				// Params: model id
				assert.Equal(t, openai.EmbeddingModel(tt.modelID), params.Model)

				// Params: input union mirrors provided values
				assert.Equal(t, tt.values, params.Input)
			})
		}
	})
}
