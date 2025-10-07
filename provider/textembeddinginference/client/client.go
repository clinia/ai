package textembeddinginference

import (
	"os"

	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

type Client struct {
	Options   []option.RequestOption
	Embedding EmbeddingService
	Reranking RerankingService
}

// DefaultClientOptions read from the environment (TEI_BASE_URL, TEI_API_KEY).
// This should be used to initialize new clients.
func DefaultClientOptions() []option.RequestOption {
	defaults := []option.RequestOption{}
	if o, ok := os.LookupEnv("TEI_BASE_URL"); ok {
		defaults = append(defaults, option.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("TEI_API_KEY"); ok {
		defaults = append(defaults, option.WithAPIKey(o))
	}
	return defaults
}

func NewClient(opts ...option.RequestOption) (r Client) {
	opts = append(DefaultClientOptions(), opts...)

	r = Client{Options: opts}

	r.Embedding = NewEmbeddingService(opts...)
	r.Reranking = NewRerankingService(opts...)
	return
}
