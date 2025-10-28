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
	values []api.MultimodalEmbeddingInput,
	opts api.TransportOptions,
) (jina.MultimodalEmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	// Map API-level inputs to Jina client inputs
	mapped := make([]jina.MultimodalEmbeddingInput, len(values))
	for i, v := range values {
		mapped[i] = jina.MultimodalEmbeddingInput{Text: v.Text, Image: v.Image}
	}

	params := jina.MultimodalEmbeddingNewParams{
		Model: jina.EmbeddingModel(modelID),
		Input: mapped,
	}

	applyProviderMultimodalMetadata(&params, opts)

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// EncodeEmbedding builds Jina params + request options from the unified API options.
func EncodeEmbedding(
	modelID string,
	values []string,
	opts api.TransportOptions,
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
func applyProviderMetadata(params *jina.TextEmbeddingNewParams, opts api.TransportOptions) {
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
func applyProviderMultimodalMetadata(params *jina.MultimodalEmbeddingNewParams, opts api.TransportOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetMultimodalEmbeddingMetadata(opts)
		if metadata != nil {
			if metadata.Task != nil && *metadata.Task != "" {
				params.Task = metadata.Task
			}
		}
	}
}
