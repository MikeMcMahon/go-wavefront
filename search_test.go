package wavefront

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

type MockSearchClient struct {
	Client
	Response  []byte
	T         *testing.T
	isDeleted bool
}

func (m MockSearchClient) Do(req *http.Request) (io.ReadCloser, error) {
	p := SearchParams{}
	b, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(b, &p)
	if err != nil {
		m.T.Fatal(err)
	}
	// check defaults
	if p.Offset != 0 || p.Limit != 100 {
		m.T.Errorf("default offset and limit, expected 0, 100; got %d, %d", p.Offset, p.Limit)
	}

	if m.isDeleted == true && req.URL.Path != "/api/v2/search/alert/deleted" {
		m.T.Errorf("deleted search path expected /api/v2/search/alert/deleted, got %s", req.URL.Path)
	}

	return io.NopCloser(bytes.NewReader(m.Response)), nil
}

func TestDefensiveCopy(t *testing.T) {
	assert := asserts.New(t)
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	response, err := os.ReadFile("./fixtures/search-alert-response.json")
	if err != nil {
		t.Fatal(err)
	}
	client := &MockSearchClient{
		Response: response,
		T:        t,
		Client: Client{
			Config:     &Config{Token: "1234-5678-9977"},
			BaseURL:    baseurl,
			httpClient: http.DefaultClient,
		},
	}
	sc := &SearchCondition{
		Key:            "id",
		Value:          "1234",
		MatchingMethod: "EXACT",
	}
	searchParams := &SearchParams{
		Conditions: []*SearchCondition{sc},
	}
	search := client.NewSearch("extlink", searchParams)

	// This step is necessary because the above call injects the real client
	// into 'search'. We want search to have our mock client instead
	search.client = client

	assert.NotSame(search.Params, searchParams)
	_, err = search.Execute()
	assert.NoError(err)
	assert.Equal(0, searchParams.Limit)
	assert.Equal(0, searchParams.Offset)
	assert.Equal(0, search.Params.Limit)
	assert.Equal(0, search.Params.Offset)
}

func TestSearch(t *testing.T) {
	assert := asserts.New(t)
	sc := &SearchCondition{
		Key:            "tags",
		Value:          "myTag",
		MatchingMethod: "EXACT",
	}

	sp := &SearchParams{
		Conditions: []*SearchCondition{sc},
	}
	response, err := os.ReadFile("./fixtures/search-alert-response.json")
	if err != nil {
		t.Fatal(err)
	}
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	s := &Search{
		Params: sp,
		Type:   "alert",
		client: &MockSearchClient{
			Response: response,
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	resp, err := s.Execute()
	if err != nil {
		t.Fatal("error executing query:", err)
	}

	raw, err := io.ReadAll(resp.RawResponse)
	if err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(raw, new(map[string]interface{})); err != nil {
		t.Error("raw response is invalid JSON", err)
	}

	// check offset of next page in paginated response
	if resp.NextOffset != 100 {
		t.Errorf("next offset, expected 100, got %d", resp.NextOffset)
	}

	// check deleted path appended
	s.Deleted = true
	((s.client).(*MockSearchClient)).isDeleted = true
	_, err = s.Execute()
	assert.NoError(err)
}
