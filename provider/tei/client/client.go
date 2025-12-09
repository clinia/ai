package tei

import (
	"os"

	"go.jetify.com/ai/provider/internal/requesterx"
)

type Client struct {
	Options   []requesterx.RequestOption
	Embedding EmbeddingService
	Ranking   RankingService
}

// DefaultClientOptions read from the environment (TEI_BASE_URL, TEI_API_KEY).
// This should be used to initialize new clients.
func DefaultClientOptions() []requesterx.RequestOption {
	defaults := []requesterx.RequestOption{}
	if o, ok := os.LookupEnv("TEI_BASE_URL"); ok {
		defaults = append(defaults, requesterx.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("TEI_API_KEY"); ok {
		defaults = append(defaults, requesterx.WithAPIKey(o))
	}
	return defaults
}

func NewClient(opts ...requesterx.RequestOption) (r Client) {
	opts = append(DefaultClientOptions(), opts...)

	r = Client{Options: opts}

	r.Embedding = NewEmbeddingService(opts...)
	r.Ranking = NewRankingService(opts...)
	return r
}
