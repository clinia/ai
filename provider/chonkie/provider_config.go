package chonkie

import (
	chonkie "go.jetify.com/ai/provider/chonkie/client"
)

type ProviderConfig struct {
	providerName string
	client       chonkie.Client
	apiKey       string
}
