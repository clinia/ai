package textembeddinginference

import (
	"context"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/textembeddinginference/client/internal/requestconfig"
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

type TextEmbeddingNewParams = embeddingNewParams[[]string]

type textEmbeddingNewParamsConcrete struct {
	TextEmbeddingNewParams
}

func (p textEmbeddingNewParamsConcrete) validate() error {
	if len(p.Inputs) == 0 {
		return fmt.Errorf("inputs: []string must be non-empty")
	}
	for i, s := range p.Inputs {
		if s == "" {
			return fmt.Errorf("inputs[%d]: empty string", i)
		}
	}
	return nil
}

// New creates an embedding vector representing the input text.
// TEI typically uses the /embed endpoint for text embeddings.
func (r *EmbeddingService) New(ctx context.Context, body TextEmbeddingNewParams, opts ...option.RequestOption) (res *CreateEmbeddingResponse, err error) {
	textEmb := textEmbeddingNewParamsConcrete{
		TextEmbeddingNewParams: body,
	}
	if err := textEmb.validate(); err != nil {
		return nil, err
	}
	opts = append(r.Options[:], opts...)
	path := "embed"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}
