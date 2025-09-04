package api

import "net/http"

// FileCallOption represent the options for generating embeddings.
type FileCallOption func(*FileCallOptions)

// WithFileCallHeaders sets HTTP headers to be sent with the request.
// Only applicable for HTTP-based providers.
func WithFileCallHeaders(headers http.Header) FileCallOption {
	return func(o *FileCallOptions) {
		o.Headers = headers
	}
}

// FileCallOptions represents the options for generating embeddings.
type FileCallOptions struct {
	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header
}
