package jina

import (
	"go.jetify.com/ai/instrumentation"
	jina "go.jetify.com/ai/provider/jina/client"
)

type ProviderConfig struct {
	providerName string
	client       jina.Client
	apiKey       string
	instrumenter instrumentation.Instrumenter
}
