package jina

import (
	"go.jetify.com/ai/api"
	jina "go.jetify.com/ai/provider/jina/client"
)

type Provider struct {
	// client is the Jina client used to make API calls.
	client jina.Client

	// name is the name of the provider, overrides the default "jina".
	name string

	// apiKey is the API key used for authentication.
	apiKey string
}

var _ api.Provider = &Provider{}

type ProviderOption func(*Provider)

func WithClient(c jina.Client) ProviderOption {
	return func(p *Provider) { p.client = c }
}

func WithName(name string) ProviderOption {
	return func(p *Provider) { p.name = name }
}

func WithAPIKey(apiKey string) ProviderOption {
	return func(p *Provider) { p.apiKey = apiKey }
}

func NewProvider(opts ...ProviderOption) api.Provider {
	p := &Provider{client: jina.NewClient()}

	for _, opt := range opts {
		opt(p)
	}

	if p.name == "" {
		p.name = "jina"
	}

	return p
}

// LanguageModel is not supported by the Jina provider.
func (p *Provider) LanguageModel(modelID string) (api.LanguageModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "LanguageModel")
}

// RankingModel is not supported by the Jina provider.
func (p *Provider) RankingModel(modelID string) (api.RankingModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "RankingModel")
}

// SparseEmbeddingModel is not supported by the Jina provider.
func (p *Provider) SparseEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.SparseEmbedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "SparseEmbeddingModel")
}
