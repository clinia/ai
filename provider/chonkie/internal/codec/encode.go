package codec

import (
	"net/http"

	"go.jetify.com/ai/provider/internal/requesterx"
)

// applyHeaders applies the provided HTTP headers to the request options.
func applyHeaders(headers http.Header) []requesterx.RequestOption {
	var reqOpts []requesterx.RequestOption
	for k, vs := range headers {
		for _, v := range vs {
			reqOpts = append(reqOpts, requesterx.WithHeaderAdd(k, v))
		}
	}
	return reqOpts
}
