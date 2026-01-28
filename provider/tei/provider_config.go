package tei

import (
	"go.jetify.com/ai/instrumentation"
	tei "go.jetify.com/ai/provider/tei/client"
)

type ProviderConfig struct {
	providerName string
	client       tei.Client
	apiKey       string
	instrumenter instrumentation.Instrumenter
}
