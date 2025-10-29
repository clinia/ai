package codec

import (
	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/tei/client"
	"go.jetify.com/ai/provider/tei/client/option"
)

// EncodeRank builds TEI rank params + request options from the unified API options.
func EncodeRank(
	query string,
	texts []string,
	opts api.TransportOptions,
) (tei.RankRequest, []option.RequestOption, []api.CallWarning, error) {
	var reqOpts []option.RequestOption
	if opts.Headers != nil {
		reqOpts = append(reqOpts, applyHeaders(opts.Headers)...)
	}

	params := tei.RankRequest{
		Query: query,
		Texts: texts,
	}

	if opts.APIKey != "" {
		reqOpts = append(reqOpts, option.WithAPIKey(opts.APIKey))
	}

	if len(opts.BaseURL) > 0 {
		reqOpts = append(reqOpts, option.WithBaseURL(opts.BaseURL))
	}

	applyRankProviderMetadata(&params, opts)

	var warnings []api.CallWarning

	return params, reqOpts, warnings, nil
}

// applyRankProviderMetadata applies metadata-specific options to the rank parameters
func applyRankProviderMetadata(params *tei.RankRequest, opts api.TransportOptions) {
	if opts.ProviderMetadata != nil {
		metadata := GetRankingMetadata(opts)
		if metadata != nil {
			if metadata.RawScores != nil {
				params.RawScores = metadata.RawScores
			}
			if metadata.ReturnText != nil {
				params.ReturnText = metadata.ReturnText
			}
			if metadata.Truncate != nil {
				params.Truncate = metadata.Truncate
			}
			if metadata.TruncationDirection != nil {
				params.TruncationDirection = metadata.TruncationDirection
			}
		}
	}
}
