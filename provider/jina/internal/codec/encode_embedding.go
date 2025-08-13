package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/client/option"
)

// EncodeEmbedding builds OpenAI params + request options from the unified API options.
func EncodeEmbedding(
	modelID string,
	values []string,
	opts api.EmbeddingOptions,
) (jina.TextEmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := jina.TextEmbeddingNewParams{
		Model: jina.EmbeddingModel(modelID),
		Input: values,
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
