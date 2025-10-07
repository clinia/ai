package textembeddinginference

import (
	"context"
	"fmt"

	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

// RerankingModel represents a TEI reranking model.
// Rerankers are used to reorder documents based on their relevance to a query.
type RerankingModel struct {
	modelID string
	pc      ProviderConfig
}

// Re-export client types for convenience
type RerankRequest = tei.RerankRequest
type RerankResponse = tei.RerankResponse
type RerankResult = tei.RerankResult

// RerankingModel creates a new TEI reranking model.
func (p *Provider) RerankingModel(modelID string) (*RerankingModel, error) {
	// Create model with provider's client
	model := &RerankingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.reranking", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}

	return model, nil
}

func (m *RerankingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *RerankingModel) SpecificationVersion() string {
	return "v1"
}

func (m *RerankingModel) ModelID() string {
	return m.modelID
}

// SupportsParallelCalls returns whether the model can handle multiple reranking calls in parallel.
func (m *RerankingModel) SupportsParallelCalls() bool {
	return true
}

// MaxDocumentsPerCall returns the maximum number of documents that can be reranked in a single call.
func (m *RerankingModel) MaxDocumentsPerCall() *int {
	max := 1000 // TODO: verify
	return &max
}

// DoRerank reranks the given documents based on their relevance to the query.
func (m *RerankingModel) DoRerank(
	ctx context.Context,
	query string,
	documents []string,
	options RerankOptions,
) (RerankResponse, error) {
	// Build the request
	request := RerankRequest{
		Query:               query,
		Texts:               documents,
		RawScores:           options.RawScores,
		ReturnText:          options.ReturnText,
		Truncate:            options.Truncate,
		TruncationDirection: options.TruncationDirection,
	}

	// Call the client
	resp, err := m.pc.client.Reranking.Rerank(ctx, request)
	if err != nil {
		return RerankResponse{}, err
	}

	return *resp, nil
}

// RerankOptions contains options for reranking requests.
type RerankOptions struct {
	// RawScores indicates whether to return raw scores instead of normalized scores
	RawScores *bool
	// ReturnText indicates whether to return the text content in the response
	ReturnText *bool
	// Truncate indicates whether to truncate inputs that are too long
	Truncate *bool
	// TruncationDirection specifies which direction to truncate from ("Left" or "Right")
	TruncationDirection *string
}
