package codec

import (
	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

// EncodeSparseEmbedding builds TEI sparse params + request options from the unified API options.
func EncodeSparseEmbedding(
	modelID string,
	values []string,
	opts api.EmbeddingOptions,
) (tei.SparseTextEmbeddingNewParams, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := tei.SparseTextEmbeddingNewParams{
		Inputs: values,
	}

	applySparseProviderMetadata(&params, opts)

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// applySparseProviderMetadata applies metadata-specific options to the sparse parameters
func applySparseProviderMetadata(params *tei.SparseTextEmbeddingNewParams, opts api.EmbeddingOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetSparseTextEmbeddingMetadata(opts)
		if metadata != nil {
			if metadata.Truncate != nil {
				params.Truncate = metadata.Truncate
			}
		}
	}
}
