package codec

import (
	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

func GetTextEmbeddingMetadata(source api.MetadataSource) *tei.TextEmbeddingNewParams {
	return api.GetMetadata[tei.TextEmbeddingNewParams]("text-embedding-inference", source)
}
