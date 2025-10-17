package codec

import (
	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
)

func GetTextEmbeddingMetadata(source api.MetadataSource) *jina.TextEmbeddingNewParams {
	return api.GetMetadata[jina.TextEmbeddingNewParams]("jina", source)
}

func GetMultimodalEmbeddingMetadata(source api.MetadataSource) *jina.MultimodalEmbeddingNewParams {
	return api.GetMetadata[jina.MultimodalEmbeddingNewParams]("jina", source)
}

// GetSegmentingMetadata retrieves per-call knobs for the Jina Segmenting.
// See jina.SegmentingNewParams for available fields.
func GetSegmentingMetadata(source api.MetadataSource) *jina.SegmentingNewParams {
	return api.GetMetadata[jina.SegmentingNewParams]("jina", source)
}
