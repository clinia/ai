package jina

// TaskKey is an optional Jina "task" hint (e.g. "text-matching", "retrieval", etc.)
type TaskKey string

// Text-only request (v3, v2, etc.). Matches Jina’s text-embedding schema.
type EmbeddingsParams struct {
	Model string   `json:"model"`          // required
	Task  *TaskKey `json:"task,omitempty"` // optional (omit for models that don't accept it)
	Input []string `json:"input"`          // required: array of strings
}

// Multimodal input element (either Text or Image, or both if the model allows).
type EmbeddingInput struct {
	Text  *string `json:"text,omitempty"`
	Image *string `json:"image,omitempty"` // URL (or base64 if your usage requires; Jina typically accepts URLs)
}

// Multimodal request (CLIP-like models). Shares the same endpoint with optional Task.
type MultimodalEmbeddingsParams struct {
	Model string           `json:"model"`          // required
	Task  *TaskKey         `json:"task,omitempty"` // optional (omit for models that don't accept it)
	Input []EmbeddingInput `json:"input"`          // required
}

// Response types based on Jina’s embeddings API shape.
type EmbeddingDatum struct {
	Object    string    `json:"object"` // "embedding"
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

type EmbeddingsResponse struct {
	Object string           `json:"object"` // "list"
	Data   []EmbeddingDatum `json:"data"`
	Model  string           `json:"model"`
	Error  *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}
