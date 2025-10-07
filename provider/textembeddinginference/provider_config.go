package textembeddinginference

import tei "go.jetify.com/ai/provider/textembeddinginference/client"

type ProviderConfig struct {
	providerName string
	client       tei.Client
	apiKey       string
}
