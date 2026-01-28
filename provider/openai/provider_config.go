package openai

import (
	"github.com/openai/openai-go/v2"
	"go.jetify.com/ai/instrumentation"
)

type ProviderConfig struct {
	providerName string
	client       openai.Client
	instrumenter instrumentation.Instrumenter
}
