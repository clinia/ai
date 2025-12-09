package jina

import (
	"os"

	"go.jetify.com/ai/provider/internal/requesterx"
	"go.jetify.com/ai/provider/jina/client/option"
)

type Client struct {
	Options    []requesterx.RequestOption
	Embeddings EmbeddingService
	Segments   SegmentingService
}

// DefaultClientOptions read from the environment (JINA_API_KEY, JINA_ORG_ID,
// JINA_PROJECT_ID, JINA_WEBHOOK_SECRET, JINA_BASE_URL). This should be used
// to initialize new clients.
func DefaultClientOptions() []requesterx.RequestOption {
	defaults := []requesterx.RequestOption{option.WithEnvironmentProduction()}
	if o, ok := os.LookupEnv("JINA_BASE_URL"); ok {
		defaults = append(defaults, requesterx.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("JINA_API_KEY"); ok {
		defaults = append(defaults, requesterx.WithAPIKey(o))
	}
	return defaults
}

func NewClient(opts ...requesterx.RequestOption) (r Client) {
	opts = append(DefaultClientOptions(), opts...)

	r = Client{Options: opts}
	r.Embeddings = NewEmbeddingService(opts...)
	r.Segments = NewSegmentingService(opts...)
	return r
}
