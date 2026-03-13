/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityH323Gatekeeper(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - initialize with defaults
	defaultPort := 1719
	mockState := &config.H323Gatekeeper{
		ID:          1,
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/1/",
		Name:        "tf-test-h323-gatekeeper",
		Description: "",
		Address:     "192.168.1.101",
		Port:        &defaultPort,
	}

	// Step 1: Create with full config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/h323_gatekeeper/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.H323GatekeeperCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
	}).Once()

	// Step 2: Update to min config
	client.On("PutJSON", mock.Anything, "configuration/v1/h323_gatekeeper/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.H323GatekeeperUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
		if args.Get(3) != nil {
			gatekeeper := args.Get(3).(*config.H323Gatekeeper)
			*gatekeeper = *mockState
		}
	}).Once()

	// Step 3: Delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/h323_gatekeeper/1/", mock.Anything).Return(nil).Maybe()

	// Step 4: Recreate with min config
	client.On("PostWithResponse", mock.Anything, "configuration/v1/h323_gatekeeper/", mock.Anything, mock.Anything).Return(&types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/h323_gatekeeper/1/",
	}, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.H323GatekeeperCreateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
	}).Once()

	// Step 5: Update to full config
	client.On("PutJSON", mock.Anything, "configuration/v1/h323_gatekeeper/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.H323GatekeeperUpdateRequest)
		mockState.Name = req.Name
		mockState.Description = req.Description
		mockState.Address = req.Address
		mockState.Port = req.Port
		if args.Get(3) != nil {
			gatekeeper := args.Get(3).(*config.H323Gatekeeper)
			*gatekeeper = *mockState
		}
	}).Once()

	// Mock Read operations (GetJSON) - used throughout all steps
	client.On("GetJSON", mock.Anything, "configuration/v1/h323_gatekeeper/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		gatekeeper := args.Get(3).(*config.H323Gatekeeper)
		*gatekeeper = *mockState
	}).Maybe()

	testInfinityH323Gatekeeper(t, client)
}

func testInfinityH323Gatekeeper(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "name", "tf-test-h323-gatekeeper"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "description", "tf-test H323 Gatekeeper Description"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "port", "1720"),
				),
			},
			{
				// Step 2: Update to min config (clear optional fields, reset to defaults)
				Config: test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "name", "tf-test-h323-gatekeeper"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "address", "192.168.1.101"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "port", "1719"),
				),
			},
			{
				// Step 3: Destroy
				Config:  test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_min"),
				Destroy: true,
			},
			{
				// Step 4: Create with min config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "name", "tf-test-h323-gatekeeper"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "address", "192.168.1.101"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "port", "1719"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_h323_gatekeeper_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "name", "tf-test-h323-gatekeeper"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "description", "tf-test H323 Gatekeeper Description"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_h323_gatekeeper.tf-test-h323-gatekeeper", "port", "1720"),
				),
			},
		},
	})
}
