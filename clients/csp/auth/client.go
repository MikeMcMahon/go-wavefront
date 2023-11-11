package auth

import (
	"crypto/tls"
	"fmt"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"net/http"
	"net/url"
	"strings"
)

// Client for CSP's authorization API
type Client struct {
	// cspAddress is the address of the csp server you want to use for auth.
	CspAddress string

	// insecure turn off TLS - only intended for testing.
	Insecure bool
}

func (a *Client) createTransport() (runtime.ClientTransport, error) {
	cfg, err := httptransport.TLSClientAuth(httptransport.TLSClientOptions{
		InsecureSkipVerify: a.Insecure,
	})
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(a.CspAddress)
	if err != nil {
		return nil, err
	}

	transport := httptransport.New(baseURL.Host, baseURL.Path, nil)
	transport.SetDebug(false)
	transport.Transport = &http.Transport{
		TLSClientConfig: cfg,
		// Not sure why this is in here.
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}

	return transport, nil
}

func (a *Client) authTokenRequest(apiToken string) *runtime.ClientOperation {
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

// GetAuthToken retrieves auth token for csp users
// When accessing and endpoint secured by CSP,
// the received `token` must be provided in the
// `Authorization` request header field as follows:
// `Authorization: Bearer {token}`
func (a *Client) GetAuthToken(apiToken string) (*OK200, error) {
	transport, err := a.createTransport()
	if err != nil {
		return nil, err
	}

	result, err := transport.Submit(a.authTokenRequest(apiToken))
	if err != nil {
		return nil, err
	}

	success, ok := result.(*OK200)

	if success != nil && strings.EqualFold(*success.Payload.TokenType, "bearer") {
		return nil, fmt.Errorf("expected a `bearer` token type, got: %s", *success.Payload.TokenType)
	}

	if ok {
		return success, nil
	}

	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response during authorization. AuthClient expected to get an error, but got: %T", result)
	panic(msg)
}
