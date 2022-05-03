package wavefront

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockMetricsPolicyClient struct {
	Client
	T *testing.T
}

type MockCrudMetricsPolicyClient struct {
	Client
	T      *testing.T
	method string
}

func (m *MockMetricsPolicyClient) Do(req *http.Request) (io.ReadCloser, error) {
	switch req.Method {
	case "GET":
		return testDo(m.T, req, "./fixtures/crud-metrics-policy-default-response.json", "GET", &MetricsPolicy{})

	case "PUT":
		return testDo(m.T, req, "./fixtures/crud-metrics-policy-response.json", "PUT", &MetricsPolicy{})

	default:
		return nil, fmt.Errorf("unimplemented METHOD %s", req.Method)
	}
}

func TestMetricsPolicy_Get(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	m := &MetricsPolicyAPI{
		client: &MockMetricsPolicyClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	id := "8bcffe68-5fcb-47fa-b935-ba7bc102b9a7"
	resp, err := m.Get()
	assert.Nil(t, err)
	assert.Equal(t, &MetricsPolicy{
		PolicyRules: []PolicyRule{{
			Accounts:    []string{},
			UserGroups:  []UserGroup{{ID: &id}},
			Roles:       []string{},
			Name:        "Allow All Metrics",
			Tags:        []string{},
			Description: "Predefined policy rule. Allows access to all metrics (timeseries, histograms, and counters) for all accounts. If this rule is removed, all accounts can access all metrics if there are no matching blocking rules.",
			Prefixes:    []string{"*"},
			TagsAnded:   false,
			AccessType:  "ALLOW",
		}},
		Customer:           "example",
		UpdaterId:          "system",
		UpdatedEpochMillis: 1603762170831,
	}, resp)
}

func TestMetricsPolicy_Post(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	m := &MetricsPolicyAPI{
		client: &MockMetricsPolicyClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	id := "8bcffe68-5fcb-47fa-b935-ba7bc102b9a7"
	id2 := "7y6ffe68-5fcb-47fa-b935-ba7bc102b9a7"
	resp, err := m.Update(&UpdateMetricsPolicyRequest{PolicyRules: []PolicyRuleRequest{{
		Accounts:     []string{},
		UserGroupIds: []string{id},
		Roles:        []string{},
		Name:         "Allow All Metrics",
		Tags:         []string{},
		Description:  "Predefined policy rule. Allows access to all metrics (timeseries, histograms, and counters) for all accounts. If this rule is removed, all accounts can access all metrics if there are no matching blocking rules.",
		Prefixes:     []string{"*"},
		TagsAnded:    true,
		AccessType:   "ALLOW",
	},
		{
			Accounts:     []string{},
			UserGroupIds: []string{id2},
			Roles:        []string{"abc123", "poi567"},
			Name:         "Allow Some Metrics",
			Tags:         []string{"Custom"},
			Description:  "Scoped filter for some.",
			Prefixes:     []string{"aa", "bb"},
			TagsAnded:    true,
			AccessType:   "DENY",
		}}})

	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp.PolicyRules))
	//assert.Equal(t,resp)
	//assert.Equal(t,resp)
	//assert.Equal(t,resp)
	//assert.Equal(t,resp)
	//assert.Equal(t,resp)
	assert.Equal(t, "example", resp.Customer)
	assert.Equal(t, "john.doe@example.com", resp.UpdaterId)
	assert.Equal(t, 2603766170831, resp.UpdatedEpochMillis)
	// TODO test object equality with pointer usergroup
	//assert.Equal(t, &MetricsPolicy{
	//	PolicyRules: []PolicyRule{{
	//		Accounts: []string{},
	//		UserGroups: []UserGroup{{
	//			ID:          &id,
	//			Name:        "Everyone",
	//			Description: "System group which contains all users",
	//		}},
	//		Roles:       []string{},
	//		Name:        "Allow All Metrics",
	//		Tags:        []string{},
	//		Description: "Predefined policy rule. Allows access to all metrics (timeseries, histograms, and counters) for all accounts. If this rule is removed, all accounts can access all metrics if there are no matching blocking rules.",
	//		Prefixes:    []string{"*"},
	//		TagsAnded:   true,
	//		AccessType:  "ALLOW",
	//	},
	//		{
	//			Accounts: []string{},
	//			UserGroups: []UserGroup{{
	//				ID:          &id2,
	//				Name:        "Some",
	//				Description: "Custom selector",
	//			}},
	//			Roles:       []string{"abc123", "poi567"},
	//			Name:        "Allow Some Metrics",
	//			Tags:        []string{"Custom"},
	//			Description: "Scoped filter for some.",
	//			Prefixes:    []string{"aa", "bb"},
	//			TagsAnded:   true,
	//			AccessType:  "DENY",
	//		}},
	//	Customer:           "example",
	//	UpdaterId:          "john.doe@example.com",
	//	UpdatedEpochMillis: 2603766170831,
	//}, resp)
}
