package api

// MultimodalEmbeddingInput represents a single multimodal input item for embedding.
// Exactly one of Text or Image should be set (non-empty).
type MultimodalEmbeddingInput struct {
	Text  *string `json:"text,omitempty"`
	Image *string `json:"image,omitempty"`
}
