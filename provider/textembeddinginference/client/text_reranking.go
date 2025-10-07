package textembeddinginference

import (
	"context"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/textembeddinginference/client/internal/requestconfig"
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

type rerankRequestConcrete struct {
	RerankRequest
}

func (p rerankRequestConcrete) validate() error {
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

// Rerank reorders the given documents based on their relevance to the query.
// Returns documents sorted by relevance score in descending order.
func (r *RerankingService) Rerank(ctx context.Context, body RerankRequest, opts ...option.RequestOption) (res *RerankResponse, err error) {
	req := rerankRequestConcrete{
		RerankRequest: body,
	}
	if err := req.validate(); err != nil {
		return nil, err
	}
	opts = append(r.Options[:], opts...)
	path := ""
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}
