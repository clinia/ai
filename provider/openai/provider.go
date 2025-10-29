package openai

import (
	"github.com/openai/openai-go/v2"
	"go.jetify.com/ai/api"
)

type Provider struct {
	// client is the OpenAI client used to make API calls.
	client openai.Client
	// name is the name of the provider, overrides the default "openai".
	name string
}

var _ api.Provider = &Provider{}

type ProviderOption func(*Provider)

func WithClient(c openai.Client) ProviderOption {
	return func(p *Provider) { p.client = c }
}

func WithName(name string) ProviderOption {
	return func(p *Provider) { p.name = name }
}

func NewProvider(opts ...ProviderOption) api.Provider {
	p := &Provider{client: openai.NewClient()}
	for _, opt := range opts {
		opt(p)
	}
	if p.name == "" {
		p.name = "openai"
	}

	return p
}

// MultimodalEmbeddingModel is not supported by OpenAI Provider at this time.
func (p *Provider) MultimodalEmbeddingModel(modelID string) (api.EmbeddingModel[api.MultimodalEmbeddingInput, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "MultiModalEmbeddingModel")
}

// SegmentingModel is not supported by OpenAI Provider at this time.
func (p *Provider) SegmentingModel(modelID string) (api.SegmentingModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "SegmentingModel")
}

// RankingModel is not supported by OpenAI Provider at this time.
func (p *Provider) RankingModel(modelID string) (api.RankingModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "RankingModel")
}

// SparseEmbeddingModel is not supported by OpenAI Provider at this time.
func (p *Provider) SparseEmbeddingModel(modelID string) (api.EmbeddingModel[string, api.SparseEmbedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "SparseEmbeddingModel")
}
