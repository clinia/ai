package chonkie

import (
	"os"

	"go.jetify.com/ai/provider/chonkie/client/option"
	"go.jetify.com/ai/provider/internal/requesterx"
)

type Client struct {
	Options    []requesterx.RequestOption
	Embeddings EmbeddingService
	Segments   SegmentingService
}

// DefaultClientOptions read from the environment (CHONKIE_API_KEY, CHONKIE_ORG_ID,
// CHONKIE_PROJECT_ID, CHONKIE_WEBHOOK_SECRET, CHONKIE_BASE_URL). This should be used
// to initialize new clients.
func DefaultClientOptions() []requesterx.RequestOption {
	defaults := []requesterx.RequestOption{option.WithEnvironmentProduction()}
	if o, ok := os.LookupEnv("CHONKIE_BASE_URL"); ok {
		defaults = append(defaults, requesterx.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("CHONKIE_API_KEY"); ok {
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
