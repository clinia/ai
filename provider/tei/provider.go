package tei

import (
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/instrumentation"
	tei "go.jetify.com/ai/provider/tei/client"
)

// Provider represents the Text Embedding Inference (TEI) provider.
type Provider struct {
	client       tei.Client
	name         string
	apiKey       string
	instrumenter instrumentation.Instrumenter
}

var _ api.Provider = &Provider{}

// ProviderOption configures the TEI provider.
type ProviderOption func(*Provider)

// WithClient sets the TEI client for the provider.
func WithClient(c tei.Client) ProviderOption {
	return func(p *Provider) {
		p.client = c
	}
}

// WithName sets the provider name for logging purposes.
func WithName(name string) ProviderOption {
	return func(p *Provider) {
		p.name = name
	}
}

// WithInstrumenter configures tracing instrumentation for provider calls.
func WithInstrumenter(instr instrumentation.Instrumenter) ProviderOption {
	return func(p *Provider) {
		if instr == nil {
			instr = instrumentation.NopInstrumenter()
		}
		p.instrumenter = instr
	}
}

// WithAPIKey sets the API key for authentication (if needed).
func WithAPIKey(apiKey string) ProviderOption {
	return func(p *Provider) {
		p.apiKey = apiKey
	}
}

// NewProvider creates a new TEI provider with the given options.
func NewProvider(opts ...ProviderOption) *Provider {
	p := &Provider{
		client:       tei.NewClient(),
		instrumenter: instrumentation.NopInstrumenter(),
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.name == "" {
		p.name = "text-embedding-inference"
	}

	return p
}

// LanguageModel is not supported by TEI provider.
func (p *Provider) LanguageModel(modelID string) (api.LanguageModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "LanguageModel")
}

// MultimodalEmbeddingModel is not supported by TEI provider.
func (p *Provider) MultimodalEmbeddingModel(modelID string) (api.EmbeddingModel[api.MultimodalEmbeddingInput, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "MultimodalEmbeddingModel")
}

// SegmentingModel is not supported by TEI provider.
func (p *Provider) SegmentingModel(modelID string) (api.SegmentingModel, error) {
	return nil, api.NewUnsupportedFunctionalityError(p.name, "SegmentingModel")
}
