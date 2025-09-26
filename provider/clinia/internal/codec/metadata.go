package codec

import (
	"github.com/clinia/models-client-go/cliniamodel/common"
	"go.jetify.com/ai/api"
)

// Metadata carries provider-specific per-call overrides extracted from ProviderMetadata.
type Metadata struct {
	// Requester allows callers/tests to inject a pre-constructed requester.
	Requester common.Requester
}

// GetMetadata extracts Clinia metadata from the given source.
func GetMetadata(source api.MetadataSource) *Metadata {
	return api.GetMetadata[Metadata]("clinia", source)
}
