package textembeddinginference

import (
	"context"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/textembeddinginference/client/internal/requestconfig"
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

// TextEmbeddingNewParams is an alias for EmbedRequest for backward compatibility
type TextEmbeddingNewParams = EmbedRequest

func validateEmbedRequest(req EmbedRequest) error {
	if req.Inputs == nil {
		return fmt.Errorf("inputs is required")
	}

	// Validate inputs based on type
	switch inputs := req.Inputs.(type) {
	case string:
		if inputs == "" {
			return fmt.Errorf("input string cannot be empty")
		}
	case []string:
		if len(inputs) == 0 {
			return fmt.Errorf("inputs array cannot be empty")
		}
		for i, s := range inputs {
			if s == "" {
				return fmt.Errorf("inputs[%d]: empty string", i)
			}
		}
	default:
		return fmt.Errorf("inputs must be string or []string")
	}

	return nil
}

// New creates an embedding vector representing the input text.
// TEI uses the /embed endpoint which returns a direct array of embedding vectors.
func (r *EmbeddingService) New(ctx context.Context, body TextEmbeddingNewParams, opts ...option.RequestOption) (res *CreateEmbeddingResponse, err error) {
	if err := validateEmbedRequest(body); err != nil {
		return nil, err
	}
	opts = append(r.Options[:], opts...)
	path := "embed" // TEI /embed endpoint as per OpenAPI spec
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}
