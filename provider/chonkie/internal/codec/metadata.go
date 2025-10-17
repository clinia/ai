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

// GetSegmentingMetadata retrieves per-call knobs for the Chonkie Segmenting.
// See chonkie.SegmentingNewParams for available fields.
func GetSegmentingMetadata(source api.MetadataSource) *chonkie.SegmentingNewParams {
	return api.GetMetadata[chonkie.SegmentingNewParams]("chonkie", source)
}
