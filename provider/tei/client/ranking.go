package tei

import (
	"go.jetify.com/ai/provider/tei/client/option"
)

// RankingService contains methods for document ranking.
type RankingService struct {
	Options []option.RequestOption
}

// RankRequest represents a request to rank documents
type RankRequest struct {
	// Query is the search query to rank documents against
	Query string `json:"query"`
	// Texts is the list of texts/documents to rank
	Texts []string `json:"texts"`
	// RawScores indicates whether to return raw scores instead of normalized scores
	RawScores *bool `json:"raw_scores,omitempty"`
	// ReturnText indicates whether to return the text content in the response
	ReturnText *bool `json:"return_text,omitempty"`
	// Truncate indicates whether to truncate inputs that are too long
	Truncate *bool `json:"truncate,omitempty"`
	// TruncationDirection specifies which direction to truncate from ("Left" or "Right")
	TruncationDirection *string `json:"truncation_direction,omitempty"`
}

// RankResponse represents the response from ranking documents
// TEI returns a direct array of Rank objects (not wrapped in a "results" field)
type RankResponse []RankResult

// RankResult represents a single ranked document (called "Rank" in TEI OpenAPI spec)
type RankResult struct {
	// Index is the original index of the text in the input list
	Index int `json:"index"`
	// Score is the relevance score (higher means more relevant)
	Score float64 `json:"score"`
	// Text is the text content (if ReturnText was true)
	Text *string `json:"text,omitempty"`
}

// NewRankingService creates a new ranking service with the given options
func NewRankingService(opts ...option.RequestOption) (r RankingService) {
	r = RankingService{}
	r.Options = opts
	return r
}
