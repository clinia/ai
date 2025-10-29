package chonkie

import (
	"os"

	"go.jetify.com/ai/provider/chonkie/client/option"
)

type Client struct {
	Options    []option.RequestOption
	Embeddings EmbeddingService
	Segments   SegmentingService
}

// DefaultClientOptions read from the environment (CHONKIE_API_KEY, CHONKIE_ORG_ID,
// CHONKIE_PROJECT_ID, CHONKIE_WEBHOOK_SECRET, CHONKIE_BASE_URL). This should be used
// to initialize new clients.
func DefaultClientOptions() []option.RequestOption {
	defaults := []option.RequestOption{option.WithEnvironmentProduction()}
	if o, ok := os.LookupEnv("CHONKIE_BASE_URL"); ok {
		defaults = append(defaults, option.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("CHONKIE_API_KEY"); ok {
		defaults = append(defaults, option.WithAPIKey(o))
	}
	return defaults
}

func NewClient(opts ...option.RequestOption) (r Client) {
	opts = append(DefaultClientOptions(), opts...)

	r = Client{Options: opts}
	r.Embeddings = NewEmbeddingService(opts...)
	r.Segments = NewSegmentingService(opts...)
	return r
}
