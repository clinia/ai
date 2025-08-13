package jina

import (
	"net/http"
)

type Provider struct {
	// client is the OpenAI client used to make API calls.
	client http.Client

	// name is the name of the provider, overrides the default "openai".
	name string

	// apiKey is the API key used for authentication.
	apiKey string
}

type ProviderOption func(*Provider)

func WithClient(c http.Client) ProviderOption {
	return func(p *Provider) { p.client = c }
}

func WithName(name string) ProviderOption {
	return func(p *Provider) { p.name = name }
}

func WithAPIKey(apiKey string) ProviderOption {
	return func(p *Provider) { p.apiKey = apiKey }
}

func NewProvider(opts ...ProviderOption) *Provider {
	p := &Provider{client: http.Client{}}

	for _, opt := range opts {
		opt(p)
	}

	if p.name == "" {
		p.name = "jina"
	}

	return p
}
