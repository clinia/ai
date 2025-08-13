package jina

import (
	jina "go.jetify.com/ai/provider/jina/client"
)

type ProviderConfig struct {
	providerName string
	client       jina.Client
	apiKey       string
}
