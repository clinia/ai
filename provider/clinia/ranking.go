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

func (p *Provider) RankingModel(modelID string) (*RankingModel, error) {
	if p.ranker == nil {
		return nil, fmt.Errorf("%s: provider ranker is nil", p.name)
	}

	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &RankingModel{
		modelID:      joinModelID(name, version),
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
	params, err := codec.EncodeRank(query, texts, opts)
	if err != nil {
		return api.RankingResponse{}, err
	}

	if m.config.ranker == nil {
		return api.RankingResponse{}, fmt.Errorf("%s: ranker is nil", m.config.providerName)
	}

	res, err := m.config.ranker.Rank(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.RankingResponse{}, err
	}

	return codec.DecodeRank(res)
}
