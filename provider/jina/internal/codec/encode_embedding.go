package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
	"go.jetify.com/ai/provider/jina/client/option"
)

// EncodeEmbedding builds Jina params + request options from the unified API options.
func EncodeMultimodalEmbedding(
	modelID string,
	values []jina.MultimodalEmbeddingInput,
	opts api.EmbeddingOptions,
) (jina.MultimodalEmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := jina.MultimodalEmbeddingNewParams{
		Model: jina.EmbeddingModel(modelID),
		Input: values,
	}

	applyProviderMultimodalMetadata(&params, opts)

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// EncodeEmbedding builds Jina params + request options from the unified API options.
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

	applyProviderMetadata(&params, opts)

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

// applyProviderMetadata applies metadata-specific options to the parameters
func applyProviderMetadata(params *jina.TextEmbeddingNewParams, opts api.EmbeddingOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetTextEmbeddingMetadata(opts)
		if metadata != nil {
			if metadata.Task != nil && *metadata.Task != "" {
				params.Task = metadata.Task
			}
		}
	}
}

// applyProviderMultimodalMetadata applies metadata-specific options to the parameters
func applyProviderMultimodalMetadata(params *jina.MultimodalEmbeddingNewParams, opts api.EmbeddingOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetMultimodalEmbeddingMetadata(opts)
		if metadata != nil {
			if metadata.Task != nil && *metadata.Task != "" {
				params.Task = metadata.Task
			}
		}
	}
}
