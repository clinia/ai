package jina

import (
	"os"

	"go.jetify.com/ai/provider/jina/client/option"
)

type Client struct {
	Options    []option.RequestOption
	Embeddings EmbeddingService
	Segmenter  SegmenterService
}

// DefaultClientOptions read from the environment (JINA_API_KEY, JINA_ORG_ID,
// JINA_PROJECT_ID, JINA_WEBHOOK_SECRET, JINA_BASE_URL). This should be used
// to initialize new clients.
func DefaultClientOptions() []option.RequestOption {
	defaults := []option.RequestOption{option.WithEnvironmentProduction()}
	if o, ok := os.LookupEnv("JINA_BASE_URL"); ok {
		defaults = append(defaults, option.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("JINA_API_KEY"); ok {
		defaults = append(defaults, option.WithAPIKey(o))
	}
	return defaults
}

func NewClient(opts ...option.RequestOption) (r Client) {
	opts = append(DefaultClientOptions(), opts...)

	r = Client{Options: opts}
	r.Embeddings = NewEmbeddingService(opts...)
	r.Segmenter = NewSegmenterService(opts...)
	return
}
