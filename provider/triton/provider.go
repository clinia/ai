package triton

import (
	"context"
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
	"go.jetify.com/ai/api"
)

// Provider wires the Clinia client options with AI SDK model interfaces.
type Provider struct {
	name          string
	clientOptions common.ClientOptions

	newEmbedder embeddingFactory
	newRanker   rankerFactory
	newChunker  chunkerFactory
	newSparse   sparseFactory
}

// Assert Provider implements the api.Provider interface
var _ api.Provider = (*Provider)(nil)

// Option configures the Provider during construction.
type Option func(*providerOptions)

type providerOptions struct {
	name          string
	clientOptions *common.ClientOptions

	newEmbedder embeddingFactory
	newRanker   rankerFactory
	newChunker  chunkerFactory
	newSparse   sparseFactory
}

// embeddingFactory defines the constructor signature for Clinia embedders.
type embeddingFactory func(context.Context, common.ClientOptions) cliniaclient.Embedder

type rankerFactory func(common.ClientOptions) cliniaclient.Ranker

type chunkerFactory func(context.Context, common.ClientOptions) cliniaclient.Chunker

type sparseFactory func(context.Context, common.ClientOptions) cliniaclient.SparseEmbedder

// WithName overrides the provider name (defaults to "clinia").
func WithName(name string) Option {
	return func(o *providerOptions) {
		o.name = name
	}
}

// WithClientOptions injects pre-built client options without a bound requester.
func WithClientOptions(opts common.ClientOptions) Option {
	cp := opts
	cp.Requester = nil
	return func(o *providerOptions) {
		o.clientOptions = &cp
	}
}

// NewProvider constructs a new Clinia provider.
func NewProvider(_ context.Context, opts ...Option) (*Provider, error) {
	options := providerOptions{name: "clinia"}
	for _, opt := range opts {
		opt(&options)
	}

	clientOpts := common.ClientOptions{}
	if options.clientOptions != nil {
		clientOpts = *options.clientOptions
		clientOpts.Requester = nil
	}

	provider := &Provider{
		name:          options.name,
		clientOptions: clientOpts,
		newEmbedder:   options.newEmbedder,
		newRanker:     options.newRanker,
		newChunker:    options.newChunker,
		newSparse:     options.newSparse,
	}

	if provider.newEmbedder == nil {
		provider.newEmbedder = cliniaclient.NewEmbedder
	}
	if provider.newRanker == nil {
		provider.newRanker = cliniaclient.NewRanker
	}
	if provider.newChunker == nil {
		provider.newChunker = cliniaclient.NewChunker
	}
	if provider.newSparse == nil {
		provider.newSparse = cliniaclient.NewSparseEmbedder
	}

	return provider, nil
}

// Name returns the provider name used for logging and metadata.
func (p *Provider) Name() string { return p.name }

// ClientOptions exposes the shared client options for advanced integrations.
func (p *Provider) ClientOptions() common.ClientOptions { return p.clientOptions }

func (p *Provider) providerNameFor(component string) string {
	return fmt.Sprintf("%s.%s", p.name, component)
}

// LanguageModel is not supported by the Clinia provider.
func (p *Provider) LanguageModel(modelID string) (api.LanguageModel, error) {
	return nil, api.NewUnsupportedFunctionalityError("language_model", "Clinia provider does not expose a language model")
}

// MultimodalEmbeddingModel is not supported by the Clinia provider.
func (p *Provider) MultimodalEmbeddingModel(modelID string) (api.EmbeddingModel[api.MultimodalEmbeddingInput, api.Embedding], error) {
	return nil, api.NewUnsupportedFunctionalityError("multimodal_embeddings", "Clinia provider does not support multimodal embeddings")
}

// withEmbeddingFactory overrides the embedder factory (used in tests).
func withEmbeddingFactory(factory embeddingFactory) Option {
	return func(o *providerOptions) {
		o.newEmbedder = factory
	}
}

func withRankerFactory(factory rankerFactory) Option {
	return func(o *providerOptions) {
		o.newRanker = factory
	}
}

func withChunkerFactory(factory chunkerFactory) Option {
	return func(o *providerOptions) {
		o.newChunker = factory
	}
}

func withSparseFactory(factory sparseFactory) Option {
	return func(o *providerOptions) {
		o.newSparse = factory
	}
}
