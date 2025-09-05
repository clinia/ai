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
