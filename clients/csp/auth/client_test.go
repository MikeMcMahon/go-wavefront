package auth

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func fixAddress(address string) string {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "https://" + address
	}
	return address + "/api/v2/"
}

func TestAccGetAuthToken(t *testing.T) {
	client := &Client{
		CspAddress: fixAddress(os.Getenv("CSP_ADDRESS")),
	}

	_, err := client.GetAuthToken(os.Getenv("WAVEFRONT_TOKEN"))
	if err != nil {
		t.Fatalf("Could not get CSP bearer token: %v", err)
	}
}

func TestAuthTokenRequest(t *testing.T) {
	// Create a test instance of the Client
	testClient := &Client{}

	// Test case: Valid API token
	apiToken := "valid_token"
	operation := testClient.authTokenRequest(apiToken)

	// Assertions
	assert.NotNil(t, operation, "CSPTokenRequest should not be nil")
	assert.Equal(t, "retrieveAuthToken", operation.ID, "CSPTokenRequest ID should match")
	assert.Equal(t, "POST", operation.Method, "HTTP method should be POST")
	assert.Equal(t, "/am/api/auth/api-tokens/authorize", operation.PathPattern, "Path pattern should match")
	assert.ElementsMatch(t, []string{"app/json", "application/json"}, operation.ProducesMediaTypes, "Produces media types should match")
	assert.ElementsMatch(t, []string{"application/json"}, operation.ConsumesMediaTypes, "Consumes media types should match")
	assert.ElementsMatch(t, []string{"https"}, operation.Schemes, "Schemes should match")

	// Check Params
	assert.NotNil(t, operation.Params, "Params should not be nil")

	// Assert that the returned Params is of type ResponseReader
	_, ok := operation.Params.(*RequestWriter)
	if !ok {
		t.Errorf("Expected the returned Params to be a RequestWriter")
	}

	// Check Body
	assert.NotNil(t, operation.Params.(*RequestWriter).Body, "Body should not be nil")
	assert.Equal(t, apiToken, *operation.Params.(*RequestWriter).Body.ApiToken, "API token in the request body should match")

	// Check Reader
	assert.NotNil(t, operation.Reader, "Reader should not be nil")

	// Assert that the returned Reader is of type ResponseReader
	_, ok = operation.Reader.(*ResponseReader)
	if !ok {
		t.Errorf("Expected the returned Reader to be a ResponseReader")
	}
}
