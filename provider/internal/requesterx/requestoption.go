package requesterx

import (
	"fmt"
	"net/http"
	"net/url"
)

func WithBaseURL(raw string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
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
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.Request.Header.Set(key, value)
		return nil
	})
}

// WithHeaderAdd returns a RequestOption that adds the header value to the associated key. It appends
// onto any existing values.
func WithHeaderAdd(key, value string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.Request.Header.Add(key, value)
		return nil
	})
}

// WithHeaderDel returns a RequestOption that deletes the header value(s) associated with the given key.
func WithHeaderDel(key string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.Request.Header.Del(key)
		return nil
	})
}

// WithHTTPClient returns a RequestOption that changes the underlying http client used to make this
// request, which by default is [http.DefaultClient].
func WithHTTPClient(c *http.Client) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		if c != nil {
			r.HTTPClient = c
		}
		return nil
	})
}

// WithAPIKey returns a RequestOption that sets the client setting "api_key".
func WithAPIKey(value string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.APIKey = value
		return r.Apply(WithHeader("authorization", fmt.Sprintf("Bearer %s", r.APIKey)))
	})
}

// WithDefaultBaseURL returns a RequestOption that sets the client's default Base URL.
// This is always overridden by setting a base URL with WithBaseURL.
// WithBaseURL should be used instead of WithDefaultBaseURL except in internal code.
func WithDefaultBaseURL(baseURL string) RequestOption {
	u, err := url.Parse(baseURL)
	return RequestOptionFunc(func(r *RequestConfig) error {
		if err != nil {
			return err
		}
		r.DefaultBaseURL = u
		return nil
	})
}

// WithUseRawBaseURL instructs the client to use the configured BaseURL as the
// full request URL without appending a path from the request.
func WithUseRawBaseURL() RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.UseRawBaseURL = true
		return nil
	})
}
