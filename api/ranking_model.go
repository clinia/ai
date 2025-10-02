package api

import "context"

// RankingModel represents a model capable of ranking a query against a set of texts.
type RankingModel interface {
	// SpecificationVersion returns which ranking model interface version is implemented.
	SpecificationVersion() string

	// ProviderName returns the name of the provider for logging purposes.
	ProviderName() string

	// ModelID returns the provider-specific model ID for logging purposes.
	ModelID() string

	// SupportsParallelCalls indicates if the model can handle multiple ranking calls in parallel.
	SupportsParallelCalls() bool

	// DoRank produces a score for each text given a query.
	DoRank(ctx context.Context, query string, texts []string, opts RankingOptions) (RankingResponse, error)
}

// RankingResponse represents the response from a ranking request.
type RankingResponse struct {
	// Scores contains one score per text in the same order as the input.
	Scores []float64

	// RequestID is an optional identifier for tracing.
	RequestID string
}
