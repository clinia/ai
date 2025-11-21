package codec

import (
	"net/http"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/requesterx"
	tei "go.jetify.com/ai/provider/tei/client"
)

// EncodeEmbedding builds TEI params + request options from the unified API options.
func EncodeEmbedding(
	modelID string,
	values []string,
	opts api.TransportOptions,
) (tei.TextEmbeddingNewParams, []requesterx.RequestOption, []api.CallWarning, error) {
	var reqOpts []requesterx.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := tei.TextEmbeddingNewParams{
		Inputs: values,
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, requesterx.WithAPIKey(opts.APIKey))
	}

	if len(opts.BaseURL) > 0 {
		reqOpts = append(reqOpts, requesterx.WithBaseURL(opts.BaseURL))
	}

	applyProviderMetadata(&params, opts)

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// applyHeaders applies the provided HTTP headers to the request options.
func applyHeaders(headers http.Header) []requesterx.RequestOption {
	var reqOpts []requesterx.RequestOption
	for k, vs := range headers {
		for _, v := range vs {
			reqOpts = append(reqOpts, requesterx.WithHeaderAdd(k, v))
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
