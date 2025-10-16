package codec

import (
	"go.jetify.com/ai/api"
	chonkie "go.jetify.com/ai/provider/chonkie/client"
)

func GetTextEmbeddingMetadata(source api.MetadataSource) *chonkie.TextEmbeddingNewParams {
	return api.GetMetadata[chonkie.TextEmbeddingNewParams]("chonkie", source)
}

func GetMultimodalEmbeddingMetadata(source api.MetadataSource) *chonkie.MultimodalEmbeddingNewParams {
	return api.GetMetadata[chonkie.MultimodalEmbeddingNewParams]("chonkie", source)
}

// GetSegmenterMetadata retrieves per-call knobs for the Chonkie Segmenter.
// See chonkie.SegmenterNewParams for available fields.
func GetSegmenterMetadata(source api.MetadataSource) *chonkie.SegmenterNewParams {
	return api.GetMetadata[chonkie.SegmenterNewParams]("chonkie", source)
}
