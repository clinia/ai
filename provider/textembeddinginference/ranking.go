package textembeddinginference

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/textembeddinginference/client"
)

// RankingModel represents a TEI ranking model.
// rankers are used to reorder documents based on their relevance to a query.
type RankingModel struct {
	modelID string
	pc      ProviderConfig
}

// Ensure rankingModel implements api.RankingModel
var _ api.RankingModel = &RankingModel{}

// Re-export client types for convenience
type RankRequest = tei.RankRequest
type RankResponse = tei.RankResponse
type RankResult = tei.RankResult

// rankingModel creates a new TEI ranking model.
func (p *Provider) RankingModel(modelID string) (*RankingModel, error) {
	// Create model with provider's client
	model := &RankingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.ranking", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
		},
	}

	return model, nil
}

func (m *RankingModel) ProviderName() string {
	return m.pc.providerName
}

func (m *RankingModel) SpecificationVersion() string {
	return "v1"
}

func (m *RankingModel) ModelID() string {
	return m.modelID
}

// SupportsParallelCalls returns whether the model can handle multiple ranking calls in parallel.
func (m *RankingModel) SupportsParallelCalls() bool {
	return true
}

// MaxDocumentsPerCall returns the maximum number of documents that can be ranked in a single call.
func (m *RankingModel) MaxDocumentsPerCall() *int {
	max := 1000 // TODO: verify
	return &max
}

// DoRank produces a score for each text given a query (implements api.RankingModel).
func (m *RankingModel) DoRank(
	ctx context.Context,
	query string,
	texts []string,
	opts api.TransportOptions,
) (api.RankingResponse, error) {
	// Build the TEI rank request
	request := RankRequest{
		Query: query,
		Texts: texts,
	}

	// Extract TEI-specific options from provider metadata if available
	if opts.ProviderMetadata != nil {
		if teiOpts, ok := opts.ProviderMetadata.Get("tei"); ok {
			if rankOpts, ok := teiOpts.(RankOptions); ok {
				request.RawScores = rankOpts.RawScores
				request.ReturnText = rankOpts.ReturnText
				request.Truncate = rankOpts.Truncate
				request.TruncationDirection = rankOpts.TruncationDirection
			}
		}
	}

	// Call the client
	resp, err := m.pc.client.Ranking.Rank(ctx, request)
	if err != nil {
		return api.RankingResponse{}, err
	}

	// Convert TEI detailed response to simple scores
	scores := make([]float64, len(texts))
	for _, result := range *resp {
		if result.Index >= 0 && result.Index < len(scores) {
			scores[result.Index] = result.Score
		}
	}

	return api.RankingResponse{
		Scores: scores,
		// RequestID could be extracted from response headers if available in the future
	}, nil
}

// RankOptions contains options for ranking requests.
type RankOptions struct {
	// RawScores indicates whether to return raw scores instead of normalized scores
	RawScores *bool
	// ReturnText indicates whether to return the text content in the response
	ReturnText *bool
	// Truncate indicates whether to truncate inputs that are too long
	Truncate *bool
	// TruncationDirection specifies which direction to truncate from ("Left" or "Right")
	TruncationDirection *string
}
