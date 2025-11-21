package tei

import (
	"context"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/internal/requesterx"
)

type rankRequestConcrete struct {
	RankRequest
}

func (p rankRequestConcrete) validate() error {
	if p.Query == "" {
		return fmt.Errorf("query is required")
	}
	if len(p.Texts) == 0 {
		return fmt.Errorf("texts: []string must be non-empty")
	}
	for i, text := range p.Texts {
		if text == "" {
			return fmt.Errorf("texts[%d]: empty string", i)
		}
	}
	// Validate truncation direction if provided
	if p.TruncationDirection != nil {
		dir := *p.TruncationDirection
		if dir != "Left" && dir != "Right" {
			return fmt.Errorf("truncation_direction must be 'Left' or 'Right'")
		}
	}
	return nil
}

// Rank reorders the given documents based on their relevance to the query.
// Returns documents sorted by relevance score in descending order.
func (r *RankingService) Rank(ctx context.Context, body RankRequest, opts ...requesterx.RequestOption) (res *RankResponse, err error) {
	req := rankRequestConcrete{
		RankRequest: body,
	}
	if err := req.validate(); err != nil {
		return nil, err
	}
	opts = append(r.Options[:], opts...)
	path := "rerank"
	err = requesterx.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}
