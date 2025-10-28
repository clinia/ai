package tei

import tei "go.jetify.com/ai/provider/tei/client"

type ProviderConfig struct {
	providerName string
	client       tei.Client
	apiKey       string
}
