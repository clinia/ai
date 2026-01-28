package chonkie

import (
	"go.jetify.com/ai/instrumentation"
	chonkie "go.jetify.com/ai/provider/chonkie/client"
)

type ProviderConfig struct {
	providerName string
	client       chonkie.Client
	apiKey       string
	instrumenter instrumentation.Instrumenter
}
