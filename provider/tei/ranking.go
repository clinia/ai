package tei

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/instrumentation"
	"go.jetify.com/ai/provider/tei/internal/codec"
)

// RankingModel represents a TEI ranking model.
// rankers are used to reorder documents based on their relevance to a query.
type RankingModel struct {
	modelID string
	pc      ProviderConfig
}

// Ensure rankingModel implements api.RankingModel
var _ api.RankingModel = &RankingModel{}

// rankingModel creates a new TEI ranking model.
func (p *Provider) RankingModel(modelID string) (api.RankingModel, error) {
	// Create model with provider's client
	model := &RankingModel{
		modelID: modelID,
		pc: ProviderConfig{
			providerName: fmt.Sprintf("%s.ranking", p.name),
			client:       p.client,
			apiKey:       p.apiKey,
			instrumenter: p.instrumenter,
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
	max := 1000 // TODO: [RET-3496] Determine actual limit
	return &max
}

// DoRank produces a score for each text given a query (implements api.RankingModel).
func (m *RankingModel) DoRank(
	ctx context.Context,
	query string,
	texts []string,
	opts api.TransportOptions,
) (resp api.RankingResponse, err error) {
	ctx, span := m.pc.instrumenter.Start(
		ctx,
		"DoRank",
		instrumentation.Attributes{
			"provider":   m.ProviderName(),
			"model":      m.modelID,
			"model_type": "ranking",
			"operation":  string(instrumentation.OperationRank),
		},
		instrumentation.ProviderSpanInfo{
			Provider:  m.ProviderName(),
			Model:     m.modelID,
			Operation: instrumentation.OperationRank,
		},
	)
	defer instrumentation.EndSpan(span, &err)

	request, reqOpts, _, err := codec.EncodeRank(query, texts, opts)
	if err != nil {
		return api.RankingResponse{}, err
	}

	// Call the client
	apiResp, err := m.pc.client.Ranking.Rank(ctx, request, reqOpts...)
	if err != nil {
		return api.RankingResponse{}, err
	}

	return codec.DecodeRank(apiResp)
}
