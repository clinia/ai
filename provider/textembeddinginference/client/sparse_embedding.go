package textembeddinginference

import (
	"context"
	"fmt"
	"net/http"

	"go.jetify.com/ai/provider/textembeddinginference/client/internal/requestconfig"
	"go.jetify.com/ai/provider/textembeddinginference/client/option"
)

// SparseTextEmbeddingNewParams is an alias for EmbedSparseRequest for backward compatibility
type SparseTextEmbeddingNewParams = EmbedSparseRequest

func validateEmbedSparseRequest(req EmbedSparseRequest) error {
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

// NewSparse creates sparse embedding vectors representing the input text.
// TEI uses the /embed_sparse endpoint which returns arrays of SparseValue objects.
func (r *EmbeddingService) NewSparse(ctx context.Context, body SparseTextEmbeddingNewParams, opts ...option.RequestOption) (res *CreateSparseEmbeddingResponse, err error) {
	if err := validateEmbedSparseRequest(body); err != nil {
		return nil, err
	}
	opts = append(r.Options[:], opts...)
	path := ""
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}
