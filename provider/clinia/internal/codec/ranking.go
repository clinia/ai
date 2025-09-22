package codec

import (
	"fmt"

	cliniaclient "github.com/clinia/models-client-go/cliniamodel"
	"go.jetify.com/ai/api"
)

// RankParams holds the resolved Clinia request.
type RankParams struct {
	ModelName    string
	ModelVersion string
	Request      cliniaclient.RankRequest
}

// EncodeRank converts the SDK call into a Clinia rank request.
func EncodeRank(modelName, modelVersion, query string, texts []string, opts api.RankingOptions) (RankParams, error) {
	if query == "" {
		return RankParams{}, fmt.Errorf("clinia/rank: query cannot be empty")
	}
	if len(texts) == 0 {
		return RankParams{}, fmt.Errorf("clinia/rank: texts cannot be empty")
	}

	req := cliniaclient.RankRequest{
		Query: query,
		Texts: texts,
	}

	return RankParams{
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Request:      req,
	}, nil
}

// DecodeRank converts the Clinia response into the SDK ranking response.
func DecodeRank(resp *cliniaclient.RankResponse) (api.RankingResponse, error) {
	if resp == nil {
		return api.RankingResponse{}, fmt.Errorf("clinia/rank: response is nil")
	}

	scores := make([]float64, len(resp.Scores))
	for i, score := range resp.Scores {
		scores[i] = float64(score)
	}

	return api.RankingResponse{
		Scores:    scores,
		RequestID: resp.ID,
	}, nil
}
