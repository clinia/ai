package option

import (
	"go.jetify.com/ai/provider/internal/requesterx"
)

// WithEnvironmentProduction returns a RequestOption that sets the current
// environment to be the "production" environment. An environment specifies which base URL
// to use by default.
func WithEnvironmentProduction() requesterx.RequestOption {
	return requesterx.WithDefaultBaseURL("https://api.jina.ai/v1/")
}
