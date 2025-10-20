package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

// EncodeEmbedding builds TEI params + request options from the unified API options.
func EncodeEmbedding(
	modelID string,
	values []string,
	opts api.TransportOptions,
) (tei.TextEmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := tei.TextEmbeddingNewParams{
		Inputs: values,
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, option.WithAPIKey(opts.APIKey))
	}

	if opts.BaseURL != nil {
		reqOpts = append(reqOpts, option.WithBaseURL(*opts.BaseURL))
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
func applyProviderMetadata(params *tei.TextEmbeddingNewParams, opts api.TransportOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetTextEmbeddingMetadata(opts)
		if metadata != nil {
			if metadata.Dimensions != nil {
				params.Dimensions = metadata.Dimensions
			}
			if metadata.Normalize != nil {
				params.Normalize = metadata.Normalize
			}
			if metadata.Truncate != nil {
				params.Truncate = metadata.Truncate
			}
			if metadata.TruncationDirection != nil {
				params.TruncationDirection = metadata.TruncationDirection
			}
			if metadata.PromptName != nil {
				params.PromptName = metadata.PromptName
			}
		}
	}
}
