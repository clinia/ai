package openai

import "github.com/openai/openai-go"

type ProviderConfig struct {
	providerName string
	client       openai.Client
}
