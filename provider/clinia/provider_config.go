package clinia

import cliniaclient "github.com/clinia/models-client-go/cliniamodel"

type ProviderConfig struct {
	providerName string
	embedder     cliniaclient.Embedder
	ranker       cliniaclient.Ranker
	chunker      cliniaclient.Chunker
}
