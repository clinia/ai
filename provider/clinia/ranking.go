package clinia

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/clinia/internal/codec"
)

type RankingModel struct {
	modelID      string
	modelName    string
	modelVersion string
	config       ProviderConfig
}

var _ api.RankingModel = (*RankingModel)(nil)

func (p *Provider) NewRankingModel(modelID string) (*RankingModel, error) {
	const defaultModelVersion = "1"

	if p.ranker == nil {
		return nil, fmt.Errorf("clinia/rank: provider ranker is nil")
	}

	name, version := splitModelIdentifier(modelID, defaultModelVersion)

	return &RankingModel{
		modelID:      modelID,
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName: p.providerNameFor("ranker"),
			ranker:       p.ranker,
		},
	}, nil
}

func (m *RankingModel) SpecificationVersion() string { return "v1" }

func (m *RankingModel) ProviderName() string { return m.config.providerName }

func (m *RankingModel) ModelID() string { return m.modelID }

func (m *RankingModel) SupportsParallelCalls() bool { return true }

func (m *RankingModel) Rank(ctx context.Context, query string, texts []string, opts api.RankingOptions) (api.RankingResponse, error) {
	params, err := codec.EncodeRank(m.modelName, m.modelVersion, query, texts, opts)
	if err != nil {
		return api.RankingResponse{}, err
	}

	if m.config.ranker == nil {
		return api.RankingResponse{}, fmt.Errorf("clinia/rank: ranker is nil")
	}

	res, err := m.config.ranker.Rank(ctx, params.ModelName, params.ModelVersion, params.Request)
	if err != nil {
		return api.RankingResponse{}, err
	}

	return codec.DecodeRank(res)
}
