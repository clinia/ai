package triton

import (
	"context"
	"fmt"

	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/triton/internal/codec"
)

type RankingModel struct {
	modelID      string
	modelName    string
	modelVersion string
	config       ProviderConfig
}

var _ api.RankingModel = (*RankingModel)(nil)

func (p *Provider) RankingModel(modelID string) (*RankingModel, error) {
	name, version, err := splitModelID(p.name, modelID)
	if err != nil {
		return nil, err
	}

	return &RankingModel{
		modelID:      modelID,
		modelName:    name,
		modelVersion: version,
		config: ProviderConfig{
			providerName:  p.providerNameFor("ranker"),
			clientOptions: p.clientOptions,
			newRanker:     p.newRanker,
		},
	}, nil
}

func (m *RankingModel) SpecificationVersion() string { return "v1" }

func (m *RankingModel) ProviderName() string { return m.config.providerName }

func (m *RankingModel) ModelID() string { return m.modelID }

func (m *RankingModel) SupportsParallelCalls() bool { return true }

func (m *RankingModel) DoRank(ctx context.Context, query string, texts []string, opts api.TransportOptions) (resp api.RankingResponse, err error) {
	params, err := codec.EncodeRank(query, texts, opts)
	if err != nil {
		return api.RankingResponse{}, err
	}

	requester, err := makeRequester(ctx, opts.BaseURL)
	if err != nil {
		return api.RankingResponse{}, err
	}
	if requester == nil {
		return api.RankingResponse{}, fmt.Errorf("%s: requester is nil", m.config.providerName)
	}

	defer func() {
		if cerr := requester.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if m.config.newRanker == nil {
		return api.RankingResponse{}, fmt.Errorf("%s: ranker factory is nil", m.config.providerName)
	}

	ranker := m.config.newRanker(m.config.clientOptionsWith(requester))
	if ranker == nil {
		return api.RankingResponse{}, fmt.Errorf("%s: ranker factory returned nil", m.config.providerName)
	}

	res, err := ranker.Rank(ctx, m.modelName, m.modelVersion, params.Request)
	if err != nil {
		return api.RankingResponse{}, err
	}

	resp, err = codec.DecodeRank(res)
	if err != nil {
		return api.RankingResponse{}, err
	}
	return resp, nil
}
