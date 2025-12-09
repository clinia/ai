package chonkie

import (
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/instrumentation"
	chonkie "go.jetify.com/ai/provider/chonkie/client"
)

type Provider struct {
	// client is the Chonkie client used to make API calls.
	client chonkie.Client

	// name is the name of the provider, overrides the default "chonkie".
	name string

	// apiKey is the API key used for authentication.
	apiKey string

	// instrumenter handles tracing spans for provider calls.
	instrumenter instrumentation.Instrumenter
}

var _ api.Provider = &Provider{}

type ProviderOption func(*Provider)

func WithClient(c chonkie.Client) ProviderOption {
	return func(p *Provider) { p.client = c }
}

func WithName(name string) ProviderOption {
	return func(p *Provider) { p.name = name }
}

func WithInstrumenter(instr instrumentation.Instrumenter) ProviderOption {
	return func(p *Provider) {
		if instr == nil {
			instr = instrumentation.NopInstrumenter()
		}
		p.instrumenter = instr
	}
}

func WithAPIKey(apiKey string) ProviderOption {
	return func(p *Provider) { p.apiKey = apiKey }
}

func NewProvider(opts ...ProviderOption) api.Provider {
	p := &Provider{
		client:       chonkie.NewClient(),
		instrumenter: instrumentation.NopInstrumenter(),
	}

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

// RankingModel is not supported by the Chonkie provider.
func (p *Provider) RankingModel(modelID string) (api.RankingModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "RankingModel")
}

// TextEmbeddingModel is not supported by the Chonkie provider.
func (p *Provider) TextEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "TextEmbeddingModel")
}

// MultimodalEmbeddingModel is not supported by the Chonkie provider.
func (p *Provider) MultimodalEmbeddingModel(modelID string) (api.EmbeddingModel[api.MultimodalEmbeddingInput, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "MultimodalEmbeddingModel")
}

// SparseEmbeddingModel is not supported by the Chonkie provider.
func (p *Provider) SparseEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.SparseEmbedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "SparseEmbeddingModel")
}
