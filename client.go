// Package wavefront provides a library for interacting with the Wavefront API,
// along with a writer for sending metrics to a Wavefront proxy.
package wavefront

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Wavefronter is an interface that a Wavefront client must satisfy
// (generally this is abstracted for easier testing)
type Wavefronter interface {
	NewRequest(method, path string, params *map[string]string, body []byte) (*http.Request, error)
	Do(req *http.Request) (io.ReadCloser, error)
}

// Config is used to hold configuration used when constructing a Client
type Config struct {
	// Address is the address of the Wavefront API, of the form example.wavefront.com
	Address string

	// Token is an authentication token that will be passed with all requests
	Token string

	// SET HTTP Proxy configuration
	HttpProxy string

	// SkipTLSVerify disables SSL certificate checking and should be used for
	// testing only
	SkipTLSVerify bool
}

// Client is used to generate API requests against the Wavefront API.
type Client struct {
	// Config is a Config object that will be used to construct requests
	Config *Config

	// BaseURL is the full URL of the Wavefront API, of the form
	// https://example.wavefront.com/api/v2
	BaseURL *url.URL

	// httpClient is the client that will be used to make requests against the API.
	httpClient *http.Client

	// debug, if set, will cause all requests to be dumped to the screen before sending.
	debug bool
}

// NewClient returns a new Wavefront client according to the given Config
func NewClient(config *Config) (*Client, error) {
	baseURL, err := url.Parse("https://" + config.Address + "/api/v2/")
	if err != nil {
		return nil, err
	}

	// need to disable http/2 as it doesn't play nicely with nginx
	// to do so we set TLSNextProto to an empty, non-nil map
	c := &Client{Config: config,
		BaseURL: baseURL,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
			},
		},
		debug: false,
	}

	// ENABLE HTTP Proxy
	if config.HttpProxy != "" {
		proxyUrl, _ := url.Parse(config.HttpProxy)
		c.httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	//For testing ONLY
	if config.SkipTLSVerify == true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.httpClient.Transport = tr
	}

	return c, nil
}

// NewRequest creates a request object to query the Wavefront API.
// Path is a relative URI that should be specified with no trailing slash,
// it will be resolved against the BaseURL of the client.
// Params should be passed as a map[string]string, these will be converted
// to query parameters.
func (c Client) NewRequest(method, path string, params *map[string]string, body []byte) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(rel)

	if params != nil {
		q := url.Query()
		for k, v := range *params {
			q.Set(k, v)
		}
		url.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.Token))
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	return req, nil
}

// Do executes a request against the Wavefront API.
// The response body is returned if the request is successful, and should
// be closed by the requester.
func (c Client) Do(req *http.Request) (io.ReadCloser, error) {

	if c.debug == true {
		d, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s\n", d)
	}

	retries := 0
	maxRetries := 4
	var buf []byte
	var err error
	if req.Body != nil {
		buf, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// reset the body since we read it already
		req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	}

	for {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// Per RFC Spec these are safe to accept as valid status codes as they all intent that the request was fulfilled
		// 200 -> OK
		// 201 -> Created
		// 202 -> Accepted
		// 203 -> Accepted but payload has been modified  via transforming proxy
		// 204 -> No Content
		if !(resp.StatusCode >= 200 && resp.StatusCode <= 204) {
			// Exponential backoff / Retry logic...
			if retries <= maxRetries {
				retries++
				// replay the buffer back into the body for retry
				if req.Body != nil {
					req.Body = ioutil.NopCloser(bytes.NewReader(buf))
				}
				time.Sleep(getSleepTime(retries))
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("server returned %s\n", resp.Status)
			}
			return nil, fmt.Errorf("server returned %s\n%s\n", resp.Status, string(body))
		}
		return resp.Body, nil
	}
}

func getSleepTime(retries int) time.Duration {
	slot := int((math.Pow(2, float64(retries)) - 1) / 2)
	rand.Seed(time.Now().UTC().UnixNano())
	var slotChoice int
	if slotChoice = 0; slot > 0 {
		slotChoice = rand.Intn(slot)
	}
	// Add some jitter, add 100ms * our random slot choice, convert to MS
	jitter := rand.Intn(50) + 50
	sleep := (time.Duration(jitter) + time.Duration(100.0*slotChoice)) * time.Millisecond
	return sleep
}

// Debug enables dumping http request objects to stdout
func (c *Client) Debug(enable bool) {
	c.debug = enable
}
