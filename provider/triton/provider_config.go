package triton

import (
	"github.com/clinia/models-client-go/cliniamodel/common"
	"go.jetify.com/ai/instrumentation"
)

type ProviderConfig struct {
	providerName  string
	clientOptions common.ClientOptions
	instrumenter  instrumentation.Instrumenter

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
