package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const SharedPolicyName = "test ingestion policy"

func AcceptanceSetup(t *testing.T) (*IngestionPolicies, *UserGroups) {
	// Configure Wavefront
	config := Config{
		Address: os.Getenv("WAVEFRONT_ADDRESS"),
		Token:   os.Getenv("WAVEFRONT_TOKEN"),
	}

	// Build Client
	client, err := NewClient(&config)

	if err != nil {
		t.Fatal(err)
	}

	// Ingestion Policy Struct
	ingestionPolicies := client.IngestionPolicies()
	userGroups := client.UserGroups()

	return ingestionPolicies, userGroups
}

func TestAccIngestionPolicy_Accounts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, _ := AcceptanceSetup(t)

	// New Policy
	policyRequest := &IngestionPolicyRequest{
		Name:        SharedPolicyName,
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "ACCOUNT",
		Accounts:    []string{"bwinter@vmware.com"},
	}

	// Create Policy
	policy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		t.Fatal(err)
	}

	err = ingestionPolicies.DeleteByID(policy.ID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccIngestionPolicy_Groups(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, userGroups := AcceptanceSetup(t)

	filter := []*SearchCondition{{
		Key:            "name",
		Value:          "Everyone",
		MatchingMethod: "EXACT",
	}}

	var groups, err = userGroups.Find(filter)

	if err != nil {
		t.Fatal(err)
	}

	// New Policy
	policyRequest := &IngestionPolicyRequest{
		Name:        SharedPolicyName,
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "GROUP",
		Groups:      []string{*groups[0].ID},
	}

	// Create Policy
	policy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		t.Fatal(err)
	}

	err = ingestionPolicies.DeleteByID(policy.ID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccIngestionPolicy_Sources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, _ := AcceptanceSetup(t)

	// New Policy
	policyRequest := &IngestionPolicyRequest{
		Name:        SharedPolicyName,
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "SOURCE",
		Sources:     []string{`source`},
	}

	// Create Policy
	policy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		t.Fatal(err)
	}

	err = ingestionPolicies.DeleteByID(policy.ID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccIngestionPolicy_Namespaces(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, _ := AcceptanceSetup(t)

	// New Policy
	policyRequest := &IngestionPolicyRequest{
		Name:        SharedPolicyName,
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "NAMESPACE",
		Namespaces:  []string{`namespace`},
	}

	// Create Policy
	policy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		t.Fatal(err)
	}

	err = ingestionPolicies.DeleteByID(policy.ID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccIngestionPolicy_Tags(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, _ := AcceptanceSetup(t)

	// New Policy
	policyRequest := &IngestionPolicyRequest{
		Name:        SharedPolicyName,
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "TAGS",
		Tags: []map[string]string{
			{"key": "user", "value": "*"},
		},
	}

	// Create Policy
	policy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		t.Fatal(err)
	}

	err = ingestionPolicies.DeleteByID(policy.ID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccIngestionPolicy_Find(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, _ := AcceptanceSetup(t)

	filter := []*SearchCondition{{
		Key:            "name",
		Value:          SharedPolicyName,
		MatchingMethod: "CONTAINS",
	}}

	var policies, err = ingestionPolicies.Find(filter)

	if err != nil {
		t.Fatal(err)
	}

	// This may be a bit flakey because it will suffer from other tests leaking.
	// Not sure an ideal solution. Would be nice to search for more than 0 elements as well.
	assert.Len(t, policies, 0)
}

func TestAccCRUDTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Acceptance Tests")
	}

	ingestionPolicies, _ := AcceptanceSetup(t)

	// New Policy
	policyRequest := &IngestionPolicyRequest{
		Name:        SharedPolicyName,
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "ACCOUNT",
		Accounts:    []string{"bwinter@vmware.com"},
	}

	// Create Policy
	ingestionPolicy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		t.Fatal(err)
	}

	ingestionPolicy, err = ingestionPolicies.GetByID(ingestionPolicy.ID)

	if err != nil {
		t.Fatal(err)
	}

	// Change the description
	ingestionPolicy.Description = "an ingestion policy updated by the Go SDK test suite"

	err = ingestionPolicies.Update(ingestionPolicy)

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	err = ingestionPolicies.DeleteByID(ingestionPolicy.ID)

	if err != nil {
		t.Fatal(err)
	}
}

type MockIngestionPoliciesClient struct {
	Client
	T *testing.T
}

type MockCrudIngestionPoliciesClient struct {
	Client
	T      *testing.T
	method string
}

func (pol MockIngestionPoliciesClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(pol.T, req, "./fixtures/search-ingestionpolicy-response.json", "POST", &SearchParams{})
}

func (pol MockCrudIngestionPoliciesClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(pol.T, req, "./fixtures/crud-ingestionpolicy-response.json", pol.method, &IngestionPolicyRequest{})
}

func TestIngestionPolicies_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	pol := &IngestionPolicies{
		client: &MockIngestionPoliciesClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				HttpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	ingestionPolicies, err := pol.Find(nil)

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 1, len(ingestionPolicies))
	assertEqual(t, "fixture-policy-id", ingestionPolicies[0].ID)
}

func TestIngestionPolicies_CreateUpdateDelete(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	pol := &IngestionPolicies{
		client: &MockCrudIngestionPoliciesClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				HttpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	pol.client.(*MockCrudIngestionPoliciesClient).method = "POST"

	policyRequest := &IngestionPolicyRequest{}
	_, err := pol.Create(policyRequest)

	if err == nil {
		t.Errorf("expected to receive error for missing fields")
	}

	policyRequest.Name = "Example"
	policyRequest.Description = "someDescription"
	policyRequest.Scope = "ACCOUNT"

	ingestionPolicy, err := pol.Create(policyRequest)
	if err != nil {
		t.Fatal(err)
	}

	pol.client.(*MockCrudIngestionPoliciesClient).method = "GET"
	var _, _ = pol.GetByID(ingestionPolicy.ID)

	policyRequest.ID = ingestionPolicy.ID
	pol.client.(*MockCrudIngestionPoliciesClient).method = "PUT"
	var _, _ = pol.rawUpdate(policyRequest)

	pol.client.(*MockCrudIngestionPoliciesClient).method = "DELETE"
	var _ = pol.DeleteByID(ingestionPolicy.ID)

	assertEqual(t, "fixture-policy-id", ingestionPolicy.ID)
}
