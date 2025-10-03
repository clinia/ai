package clinia

import (
	"github.com/clinia/models-client-go/cliniamodel/common"
)

type ProviderConfig struct {
	providerName  string
	clientOptions common.ClientOptions

	newEmbedder embeddingFactory
	newRanker   rankerFactory
	newChunker  chunkerFactory
	newSparse   sparseFactory
}

func (c ProviderConfig) clientOptionsWith(requester common.Requester) common.ClientOptions {
	opts := c.clientOptions
	opts.Requester = requester
	return opts
}
