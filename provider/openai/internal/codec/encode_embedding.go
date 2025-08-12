package codec

import (
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.jetify.com/ai/api"
)

// EncodeEmbedding builds OpenAI params + request options from the unified API options.
func EncodeEmbedding(
	modelID string,
	values []string,
	opts []api.EmbeddingOption,
) (openai.EmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	eo := api.EmbeddingOptions{}
	for _, opt := range opts {
		opt(&eo)
	}

	var reqOpts []option.RequestOption
	if eo.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(eo.Headers)...)
	}

	if eo.BaseURL != nil {
		reqOpts = append(reqOpts, option.WithBaseURL(*eo.BaseURL))
	}

	params := openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(modelID),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: values,
		},
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	}

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// applyHeaders applies the provided HTTP headers to the request options.
func applyHeaders(headers http.Header) []option.RequestOption {
	var reqOpts []option.RequestOption
	for k, vs := range headers {
		for _, v := range vs {
			reqOpts = append(reqOpts, option.WithHeaderAdd(k, v))
		}
	}
	return reqOpts
}
