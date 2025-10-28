package codec

import (
	"go.jetify.com/ai/api"
	tei "go.jetify.com/ai/provider/tei/client"
)

// DecodeRank maps the TEI ranking API response to the unified api.RankingResponse.
// TEI returns an array of RankResult objects with index and score fields.
func DecodeRank(resp *tei.RankResponse) (api.RankingResponse, error) {
	if resp == nil {
		return api.RankingResponse{}, api.NewEmptyResponseBodyError("response from TEI ranking API is nil")
	}

	// TEI returns []RankResult
	rankResults := *resp

	// We need to determine the number of original texts to create the scores array
	// The scores array should match the order of the input texts
	maxIndex := -1
	for _, result := range rankResults {
		if result.Index > maxIndex {
			maxIndex = result.Index
		}
	}

	if maxIndex < 0 {
		return api.RankingResponse{
			Scores: []float64{},
		}, nil
	}

	// Create scores array with the right size
	scores := make([]float64, maxIndex+1)

	// Fill in the scores at their respective indices
	for _, result := range rankResults {
		if result.Index >= 0 && result.Index < len(scores) {
			scores[result.Index] = result.Score
		}
	}

	return api.RankingResponse{
		Scores: scores,
	}, nil
}
