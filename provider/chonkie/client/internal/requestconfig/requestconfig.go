package requestconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestOption is applied when preparing an HTTP request (headers, base URL, etc.).
type RequestOption interface {
	Apply(*RequestConfig) error
}

type RequestOptionFunc func(*RequestConfig) error

func (s RequestOptionFunc) Apply(r *RequestConfig) error { return s(r) }

// RequestConfig holds reusable request settings for the Chonkie client.
type RequestConfig struct {
	BaseURL        *url.URL
	DefaultBaseURL *url.URL // optional fallback if BaseURL is not set
	Request        *http.Request
	HTTPClient     *http.Client
	APIKey         string
	// UseRawBaseURL instructs the executor to use the BaseURL as the full
	// request URL without resolving the request path against it.
	UseRawBaseURL bool
	// If ResponseBodyInto not nil, then we will attempt to deserialize into
	// ResponseBodyInto. If Destination is a []byte, then it will return the body as
	// is.
	ResponseBodyInto any
}

// NewRequestConfig returns a minimal config with sensible defaults.
func NewRequestConfig(
	ctx context.Context,
	method, urlStr string,
	body any,
	dst any,
	opts ...RequestOption,
) (*RequestConfig, error) {
	var reader io.Reader
	if body != nil {
		content, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(content)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	cfg := &RequestConfig{
		Request:          req,
		HTTPClient:       http.DefaultClient,
		ResponseBodyInto: dst,
	}

	if err := cfg.Apply(opts...); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ExecuteNewRequest(
	ctx context.Context,
	method, urlStr string,
	body any,
	dst any,
	opts ...RequestOption,
) error {
	cfg, err := NewRequestConfig(ctx, method, urlStr, body, dst, opts...)
	if err != nil {
		return err
	}
	return cfg.Execute()
}

func (cfg *RequestConfig) Execute() error {
	if cfg.BaseURL == nil {
		if cfg.DefaultBaseURL != nil {
			cfg.BaseURL = cfg.DefaultBaseURL
		} else {
			return fmt.Errorf("requestconfig: base url is not set")
		}
	}

	// Build final URL
	if cfg.UseRawBaseURL {
		// Ignore the request path and use BaseURL as-is
		cfg.Request.URL = cfg.BaseURL
	} else {
		// Resolve the request path relative to BaseURL
		u := cfg.BaseURL.ResolveReference(cfg.Request.URL)
		cfg.Request.URL = u
	}

	resp, err := cfg.HTTPClient.Do(cfg.Request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chonkie API error: status=%d body=%s", resp.StatusCode, string(body))
	}

	if cfg.ResponseBodyInto == nil {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return json.NewDecoder(resp.Body).Decode(cfg.ResponseBodyInto)
	}

	// Non-JSON: assume []byte or string
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	switch dst := cfg.ResponseBodyInto.(type) {
	case *[]byte:
		*dst = data
	case *string:
		*dst = string(data)
	default:
		return fmt.Errorf("unsupported destination type %T", dst)
	}
	return nil
}

// Apply applies each option in order.
func (cfg *RequestConfig) Apply(opts ...RequestOption) error {
	for _, opt := range opts {
		if err := opt.Apply(cfg); err != nil {
			return err
		}
	}
	return nil
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
