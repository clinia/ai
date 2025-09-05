package jina

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/jina/client/internal/requestconfig"
	"go.jetify.com/ai/provider/jina/client/option"
)

type MultimodalEmbeddingNewParams = embeddingNewParams[[]MultimodalEmbeddingInput]

type multimodalNewParamsConcrete struct {
	MultimodalEmbeddingNewParams
}

// MultimodalEmbeddingInput matches Jina's multimodal item shape (one of text or image).
type MultimodalEmbeddingInput struct {
	Text  *string `json:"text,omitempty"`
	Image *string `json:"image,omitempty"`
}

func (p multimodalNewParamsConcrete) validate() error {
	if p.Model == "" {
		return fmt.Errorf("model is required")
	}
	if len(p.Input) == 0 {
		return fmt.Errorf("input: []MultiModalEmbeddingInput must be non-empty")
	}
	for i, mm := range p.Input {
		if err := mm.validate(); err != nil {
			return fmt.Errorf("input[%d]: %w", i, err)
		}
	}
	return nil
}

func (it MultimodalEmbeddingInput) validate() error {
	hasText := it.Text != nil && *it.Text != ""
	hasImage := it.Image != nil && *it.Image != ""
	switch {
	case hasText && hasImage:
		return errors.New("MultiModalEmbeddingInput: exactly one of text or image must be set (not both)")
	case !hasText && !hasImage:
		return errors.New("MultiModalEmbeddingInput: one of text or image must be set")
	default:
		return nil
	}
}

// Creates an embedding vector representing the multi-modal input.
func (r *EmbeddingService) NewMultiModal(ctx context.Context, body MultimodalEmbeddingNewParams, opts ...option.RequestOption) (res *CreateEmbeddingResponse, err error) {
	textEmb := multimodalNewParamsConcrete{
		MultimodalEmbeddingNewParams: body,
	}
	if err := textEmb.validate(); err != nil {
		return nil, err
	}

	opts = append(r.Options[:], opts...)
	path := "embeddings"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}
