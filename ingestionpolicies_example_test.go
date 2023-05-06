package wavefront_test

import (
	"fmt"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleIngestionPolicies() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}

	client, err := wavefront.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	ingestionPolicies := client.IngestionPolicies()

	policyRequest := &wavefront.IngestionPolicyRequest{
		Name:        "test ingestion policy",
		Description: "an ingestion policy created by the Go SDK test suite",
		Scope:       "ACCOUNT",
		Accounts:    []string{"user@example.com"},
	}

	ingestionPolicy, err := ingestionPolicies.Create(policyRequest)

	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the policy
	fmt.Println("policy ID is", ingestionPolicy.ID)

	ingestionPolicy, err = ingestionPolicies.GetByID(ingestionPolicy.ID)

	if err != nil {
		log.Fatal(err)
	}

	// Change the description
	ingestionPolicy.Description = "an ingestion policy updated by the Go SDK test suite"

	err = ingestionPolicies.Update(ingestionPolicy)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	err = ingestionPolicies.DeleteByID(ingestionPolicy.ID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("policy deleted")
}
