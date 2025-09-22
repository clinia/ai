package clinia

import (
	"context"
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
)

// Provider wires the Clinia requester with AI SDK model interfaces.
type Provider struct {
	name          string
	clientOptions common.ClientOptions
	embedder      cliniaclient.Embedder
	ranker        cliniaclient.Ranker
	chunker       cliniaclient.Chunker
}

// Option configures the Provider during construction.
type Option func(*providerOptions)

type providerOptions struct {
	name          string
	requester     common.Requester
	clientOptions *common.ClientOptions
}

// WithName overrides the provider name (defaults to "clinia").
func WithName(name string) Option {
	return func(o *providerOptions) {
		o.name = name
	}
}

// WithRequester supplies an explicit requester to be reused by all models.
func WithRequester(r common.Requester) Option {
	return func(o *providerOptions) {
		o.requester = r
	}
}

// WithClientOptions injects pre-built client options.
func WithClientOptions(opts common.ClientOptions) Option {
	cp := opts
	return func(o *providerOptions) {
		o.clientOptions = &cp
	}
}

// NewProvider constructs a new Clinia provider. A requester must be supplied through
// WithRequester or WithClientOptions.
func NewProvider(ctx context.Context, opts ...Option) (*Provider, error) {
	options := providerOptions{name: "clinia"}
	for _, opt := range opts {
		opt(&options)
	}

	clientOpts := common.ClientOptions{}
	if options.clientOptions != nil {
		clientOpts = *options.clientOptions
	}

	if clientOpts.Requester == nil {
		clientOpts.Requester = options.requester
	}

	if clientOpts.Requester == nil {
		return nil, fmt.Errorf("clinia/provider: requester is required")
	}
	embedder := cliniaclient.NewEmbedder(ctx, clientOpts)
	ranker := cliniaclient.NewRanker(clientOpts)
	chunker := cliniaclient.NewChunker(ctx, clientOpts)

	return &Provider{
		name:          options.name,
		clientOptions: clientOpts,
		embedder:      embedder,
		ranker:        ranker,
		chunker:       chunker,
	}, nil
}

// Name returns the provider name used for logging and metadata.
func (p *Provider) Name() string { return p.name }

// ClientOptions exposes the shared client options for advanced integrations.
func (p *Provider) ClientOptions() common.ClientOptions { return p.clientOptions }

// Embedder exposes the underlying Clinia embedder implementation.
func (p *Provider) Embedder() cliniaclient.Embedder { return p.embedder }

// Ranker exposes the underlying Clinia ranker implementation.
func (p *Provider) Ranker() cliniaclient.Ranker { return p.ranker }

// Chunker exposes the underlying Clinia chunker implementation.
func (p *Provider) Chunker() cliniaclient.Chunker { return p.chunker }

func (p *Provider) providerNameFor(component string) string {
	return fmt.Sprintf("%s.%s", p.name, component)
}
