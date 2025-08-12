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
	var eo api.EmbeddingOptions
	for _, opt := range opts {
		opt(&eo)
	}

	var reqOpts []option.RequestOption
	if eo.Headers != nil {
		applyHeaders := func(h http.Header) {
			for k, vs := range h {
				for _, v := range vs {
					reqOpts = append(reqOpts, option.WithHeader(k, v))
				}
			}
		}
		applyHeaders(eo.Headers)
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
