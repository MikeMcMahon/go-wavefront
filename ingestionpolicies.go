package wavefront

import (
	"fmt"
)

const baseIngestionPolicyPath = "/api/v2/usage/ingestionpolicy"

type IngestionPolicyAccount struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type IngestionPolicyGroup struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

type IngestionPolicyTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type IngestionPolicyResponse struct {
	ID          string                   `json:"id,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	Scope       string                   `json:"scope,omitempty"`
	Accounts    []IngestionPolicyAccount `json:"accounts,omitempty"`
	Groups      []IngestionPolicyGroup   `json:"groups,omitempty"`
	Sources     []string                 `json:"sources,omitempty"`
	Namespaces  []string                 `json:"namespaces,omitempty"`
	Tags        []IngestionPolicyTag     `json:"pointTags,omitempty"`
}

type IngestionPolicyRequest struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Scope       string   `json:"scope,omitempty"`
	Accounts    []string `json:"accounts,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Sources     []string `json:"sources,omitempty"`
	Namespaces  []string `json:"namespaces,omitempty"`
	// This is a slice of maps because you may have different values with the same key
	// and also there is potentially an AND or OR operation assigned to them.
	Tags []map[string]string `json:"pointTags,omitempty"`
}

type IngestionPolicies struct {
	client Wavefronter
}

func (c *Client) IngestionPolicies() *IngestionPolicies {
	return &IngestionPolicies{client: c}
}

func (p *IngestionPolicies) Find(conditions []*SearchCondition) (results []*IngestionPolicyResponse, err error) {

	err = doSearch(conditions, "ingestionpolicy", p.client, &results)
	return // results, err
}

func (p *IngestionPolicies) GetByID(policyID string) (*IngestionPolicyResponse, error) {
	if policyID == "" {
		return nil, fmt.Errorf("id must be specified")
	}

	ingestionPolicy := IngestionPolicyResponse{}

	err := doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseIngestionPolicyPath, policyID),
		p.client,
		doResponse(&ingestionPolicy),
	)

	if err != nil {
		return nil, err
	}

	return &ingestionPolicy, nil
}

func (p *IngestionPolicies) Create(policy *IngestionPolicyRequest) (*IngestionPolicyResponse, error) {
	if policy.Name == "" {
		return nil, fmt.Errorf("ingestion policy name must be specified")
	}

	ingestionPolicy := IngestionPolicyResponse{}

	err := doRest(
		"POST",
		baseIngestionPolicyPath,
		p.client,
		doPayload(policy),
		doResponse(&ingestionPolicy),
	)

	if err != nil {
		return nil, err
	}

	return &ingestionPolicy, nil
}

func (p *IngestionPolicies) Update(policy *IngestionPolicyResponse) error {

	var accounts []string
	for _, v := range policy.Accounts {
		accounts = append(accounts, v.ID)
	}

	var groups []string
	for _, v := range policy.Groups {
		groups = append(groups, v.ID)
	}

	var pointTags []map[string]string
	for _, v := range policy.Tags {
		pointTags = append(pointTags, map[string]string{"Key": v.Key, "Value": v.Value})
	}

	policyRequest := &IngestionPolicyRequest{
		ID:          policy.ID,
		Name:        policy.Name,
		Description: policy.Description,
		Scope:       policy.Scope,
		Accounts:    accounts,
		Groups:      groups,
		Sources:     policy.Sources,
		Namespaces:  policy.Namespaces,
		Tags:        pointTags,
	}

	updatedPolicy, err := p.rawUpdate(policyRequest)

	if err != nil {
		return err
	}

	*policy = *updatedPolicy

	return nil
}

func (p *IngestionPolicies) rawUpdate(policy *IngestionPolicyRequest) (*IngestionPolicyResponse, error) {
	if policy.ID == "" {
		return nil, fmt.Errorf("id must be specified")
	}

	ingestionPolicy := IngestionPolicyResponse{}

	err := doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseIngestionPolicyPath, policy.ID),
		p.client,
		doPayload(policy),
		doResponse(&ingestionPolicy),
	)

	if err != nil {
		return nil, err
	}

	return &ingestionPolicy, nil
}

func (p *IngestionPolicies) DeleteByID(policyID string) error {
	if policyID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseIngestionPolicyPath, policyID),
		p.client,
	)
}
