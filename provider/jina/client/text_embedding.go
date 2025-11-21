package jina

import (
	"context"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/internal/requesterx"
)

type TextEmbeddingNewParams = embeddingNewParams[[]string]

type textEmbeddingNewParamsConcrete struct {
	TextEmbeddingNewParams
}

func (p textEmbeddingNewParamsConcrete) validate() error {
	if p.Model == "" {
		return fmt.Errorf("model is required")
	}
	if len(p.Input) == 0 {
		return fmt.Errorf("input: []string must be non-empty")
	}
	for i, s := range p.Input {
		if s == "" {
			return fmt.Errorf("input[%d]: empty string", i)
		}
	}
	return nil
}

// Creates an embedding vector representing the input text.
func (r *EmbeddingService) New(ctx context.Context, body TextEmbeddingNewParams, opts ...requesterx.RequestOption) (res *CreateEmbeddingResponse, err error) {
	textEmb := textEmbeddingNewParamsConcrete{
		TextEmbeddingNewParams: body,
	}
	if err := textEmb.validate(); err != nil {
		return nil, err
	}
	opts = append(r.Options[:], opts...)
	path := "embeddings"
	err = requesterx.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	if err != nil {
		return nil, fmt.Errorf("jina embedding new: %w", err)
	}
	return res, nil
}
