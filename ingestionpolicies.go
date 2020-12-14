package wavefront

import (
	"fmt"
)

const basePolicyPath = "/api/v2/usage/ingestionpolicy"

type IngestionPolicy struct {
	ID                   *string `json:"id,omitempty"`
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	lastUpdatedMs        int     `json:"lastUpdatedMs"`
	lastUpdatedAccountId int     `json:"lastUpdatedAccountId"`
	UserAccountCount     int     `json:"userAccountCount"`
	ServiceAccountCount  int     `json:"serviceAccountCount"`
}

type IngestionPolicies struct {
	client Wavefronter
}

func (c *Client) IngestionPolicies() *IngestionPolicies {
	return &IngestionPolicies{client: c}
}

func (e IngestionPolicies) Find(conditions []*SearchCondition) (
	results []*IngestionPolicy, err error) {
	err = doSearch(conditions, "ingestionpolicy", e.client, &results)
	return
}

func (e IngestionPolicies) Get(policy *IngestionPolicy) error {
	if policy.ID == nil || *policy.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", basePolicyPath, *policy.ID),
		e.client,
		doResponse(policy))
}

func (e IngestionPolicies) Create(policy *IngestionPolicy) error {
	if policy.Name == "" {
		return fmt.Errorf("ingestion policy name must be specified")
	}
	return doRest(
		"POST",
		basePolicyPath,
		e.client,
		doPayload(policy),
		doResponse(policy))
}

func (e IngestionPolicies) Update(policy *IngestionPolicy) error {
	if policy.ID == nil || *policy.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", basePolicyPath, *policy.ID),
		e.client,
		doPayload(policy),
		doResponse(policy))
}

func (e IngestionPolicies) Delete(policy *IngestionPolicy) error {
	if policy.ID == nil || *policy.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", basePolicyPath, *policy.ID),
		e.client)
	if err != nil {
		return err
	}

	// Clear out the id to prevent re-submission
	empty := ""
	policy.ID = &empty
	return nil
}
