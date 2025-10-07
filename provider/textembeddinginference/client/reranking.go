package textembeddinginference

import (
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

// RerankingService contains methods for document reranking.
type RerankingService struct {
	Options []option.RequestOption
}

// RerankRequest represents a request to rerank documents
type RerankRequest struct {
	// Query is the search query to rank documents against
	Query string `json:"query"`
	// Texts is the list of texts/documents to rerank
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

// RerankResponse represents the response from reranking documents
// TEI returns a direct array of Rank objects (not wrapped in a "results" field)
type RerankResponse []RerankResult

// RerankResult represents a single reranked document (called "Rank" in TEI OpenAPI spec)
type RerankResult struct {
	// Index is the original index of the text in the input list
	Index int `json:"index"`
	// Score is the relevance score (higher means more relevant)
	Score float64 `json:"score"`
	// Text is the text content (if ReturnText was true)
	Text *string `json:"text,omitempty"`
}

// NewRerankingService creates a new reranking service with the given options
func NewRerankingService(opts ...option.RequestOption) (r RerankingService) {
	r = RerankingService{}
	r.Options = opts
	return r
}
