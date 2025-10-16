package chonkie

import (
	"go.jetify.com/ai/api"
	chonkie "go.jetify.com/ai/provider/chonkie/client"
)

type Provider struct {
	// client is the Chonkie client used to make API calls.
	client chonkie.Client

	// name is the name of the provider, overrides the default "chonkie".
	name string

	// apiKey is the API key used for authentication.
	apiKey string
}

var _ api.Provider = &Provider{}

type ProviderOption func(*Provider)

func WithClient(c chonkie.Client) ProviderOption {
	return func(p *Provider) { p.client = c }
}

func WithName(name string) ProviderOption {
	return func(p *Provider) { p.name = name }
}

func WithAPIKey(apiKey string) ProviderOption {
	return func(p *Provider) { p.apiKey = apiKey }
}

func NewProvider(opts ...ProviderOption) *Provider {
	p := &Provider{client: chonkie.NewClient()}

	for _, opt := range opts {
		opt(p)
	}

	if p.name == "" {
		p.name = "chonkie"
	}

	return p
}

// LanguageModel is not supported by the Chonkie provider.
func (p *Provider) LanguageModel(modelID string) (api.LanguageModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "LanguageModel")
}
