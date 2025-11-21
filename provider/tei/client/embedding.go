package tei

import (
	"go.jetify.com/ai/provider/tei/client/option"
)

// EmbeddingService contains methods and other services that help with interacting
// with the embedding API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEmbeddingService] method instead.
type EmbeddingService struct {
	Options []option.RequestOption
}

// CreateEmbeddingResponse represents the response from the TEI /embed endpoint
// TEI returns embeddings as a direct array of arrays of floats: [[0.0, 1.0, 2.0], [0.1, 1.1, 2.1]]
type CreateEmbeddingResponse [][]float64

// CreateSparseEmbeddingResponse represents the response from the TEI /embed_sparse endpoint
// TEI returns sparse embeddings as an array of arrays of SparseValue objects
type CreateSparseEmbeddingResponse []SparseValue

// SparseValue represents a single non-zero value in a sparse embedding
// This matches the TEI API SparseValue schema
type SparseValue map[string]float64

// TokenWeight represents a single token's index and its corresponding weight/value
type TokenWeight struct {
	Token  string  `json:"token"`
	Weight float64 `json:"weight"`
}

// EmbeddingModel represents a TEI model identifier
type EmbeddingModel = string

// EmbedRequest represents the request parameters for the TEI /embed endpoint
// This matches the TEI OpenAPI EmbedRequest schema
type EmbedRequest struct {
	// Inputs can be a single string or array of strings to embed
	Inputs interface{} `json:"inputs"`
	// Dimensions is the number of dimensions that the output embeddings should have (optional)
	Dimensions *int `json:"dimensions,omitempty"`
	// Normalize indicates whether to normalize the embeddings (default: true)
	Normalize *bool `json:"normalize,omitempty"`
	// Truncate indicates whether to truncate inputs that exceed model limits (default: false)
	Truncate *bool `json:"truncate,omitempty"`
	// TruncationDirection specifies truncation direction: "Left" or "Right" (default: "Right")
	TruncationDirection *string `json:"truncation_direction,omitempty"`
	// PromptName is the name of the prompt to use for encoding (optional)
	PromptName *string `json:"prompt_name,omitempty"`
}

// EmbedSparseRequest represents the request parameters for the TEI /embed_sparse endpoint
// This matches the TEI OpenAPI EmbedSparseRequest schema
type EmbedSparseRequest struct {
	// Inputs can be a single string or array of strings to embed
	Inputs interface{} `json:"inputs"`
	// Truncate indicates whether to truncate inputs that exceed model limits (default: false)
	Truncate *bool `json:"truncate,omitempty"`
	// TruncationDirection specifies truncation direction: "Left" or "Right" (default: "Right")
	TruncationDirection *string `json:"truncation_direction,omitempty"`
	// PromptName is the name of the prompt to use for encoding (optional)
	PromptName *string `json:"prompt_name,omitempty"`
}

// NewEmbeddingService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEmbeddingService(opts ...option.RequestOption) (r EmbeddingService) {
	r = EmbeddingService{}
	r.Options = opts
	return r
}
