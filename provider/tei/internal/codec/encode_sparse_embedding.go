package codec

import (
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/requesterx"
	tei "go.jetify.com/ai/provider/tei/client"
)

// EncodeSparseEmbedding builds TEI sparse params + request options from the unified API options.
func EncodeSparseEmbedding(
	modelID string,
	values []string,
	opts api.TransportOptions,
) (tei.SparseTextEmbeddingNewParams, []requesterx.RequestOption, []api.CallWarning, error) {
	var reqOpts []requesterx.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := tei.SparseTextEmbeddingNewParams{
		Inputs: values,
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, requesterx.WithAPIKey(opts.APIKey))
	}

	if len(opts.BaseURL) > 0 {
		reqOpts = append(reqOpts, requesterx.WithBaseURL(opts.BaseURL))
	}

	applySparseProviderMetadata(&params, opts)

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// applySparseProviderMetadata applies metadata-specific options to the sparse parameters
func applySparseProviderMetadata(params *tei.SparseTextEmbeddingNewParams, opts api.TransportOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetSparseTextEmbeddingMetadata(opts)
		if metadata != nil {
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
