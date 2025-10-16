package option

import (
	"fmt"
	"net/http"
	"net/url"

	"go.jetify.com/ai/provider/chonkie/client/internal/requestconfig"
)

// RequestOption is an option for the requests made by the chonkie API Client
// which can be supplied to clients, services, and methods.
type RequestOption = requestconfig.RequestOption

func WithBaseURL(raw string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		r.BaseURL = u
		return nil
	})
}

// WithHeader returns a RequestOption that sets the header value to the associated key. It overwrites
// any value if there was one already present.
func WithHeader(key, value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Request.Header.Set(key, value)
		return nil
	})
}

// WithHeaderAdd returns a RequestOption that adds the header value to the associated key. It appends
// onto any existing values.
func WithHeaderAdd(key, value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Request.Header.Add(key, value)
		return nil
	})
}

// WithHeaderDel returns a RequestOption that deletes the header value(s) associated with the given key.
func WithHeaderDel(key string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Request.Header.Del(key)
		return nil
	})
}

// WithHTTPClient returns a RequestOption that changes the underlying http client used to make this
// request, which by default is [http.DefaultClient].
func WithHTTPClient(c *http.Client) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		if c != nil {
			r.HTTPClient = c
		}
		return nil
	})
}

// WithAPIKey returns a RequestOption that sets the client setting "api_key".
func WithAPIKey(value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.APIKey = value
		return r.Apply(WithHeader("authorization", fmt.Sprintf("Bearer %s", r.APIKey)))
	})
}

// WithUseRawBaseURL instructs the client to use the configured BaseURL as the
// full request URL without appending a path from the request.
func WithUseRawBaseURL() RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.UseRawBaseURL = true
		return nil
	})
}

// WithEnvironmentProduction returns a RequestOption that sets the current
// environment to be the "production" environment. An environment specifies which base URL
// to use by default.
func WithEnvironmentProduction() RequestOption {
	return requestconfig.WithDefaultBaseURL("https://api.chonkie.ai/v1/")
}
