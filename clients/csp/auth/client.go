package auth

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"net/url"
)

// HTTPTransport creates a generic http transport using a *url.URL and an optional "insecure" bool.
func HTTPTransport(baseURL *url.URL, insecure bool) (runtime.ClientTransport, error) {
	return nil, nil
}

func CSPTokenRequest(apiToken string) *runtime.ClientOperation {
	return &runtime.ClientOperation{
		ID:                 "retrieveAuthToken",
		Method:             "POST",
		PathPattern:        "/am/api/auth/api-tokens/authorize",
		ProducesMediaTypes: []string{"app/json", "application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},

		Params: &RequestWriter{
			timeout: httptransport.DefaultTimeout,
			Body:    &RequestBody{ApiToken: &apiToken},
		},

		Reader: &ResponseReader{
			formats: strfmt.Default,
		},
	}
}
