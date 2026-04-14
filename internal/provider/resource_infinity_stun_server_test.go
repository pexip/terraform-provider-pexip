/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityStunServer(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.STUNServer{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/stun_server/123/",
		Name:        "tf-test-stun-server-full",
		Description: "tf-test STUN server description",
		Address:     "stun-full.example.com",
		Port:        5349,
	}

	// Mock the CreateSTUNServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/stun_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/stun_server/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.STUNServerCreateRequest)
		// Update mock state based on create request
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		mockState.Address = createReq.Address
		mockState.Port = createReq.Port
	}).Maybe()

	// Mock the GetSTUNServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/stun_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		stunServer := args.Get(3).(*config.STUNServer)
		*stunServer = *mockState
	}).Maybe()

	// Mock the UpdateSTUNServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/stun_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.STUNServerUpdateRequest)
		stunServer := args.Get(3).(*config.STUNServer)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		// Description always sent now (without omitempty)
		mockState.Description = updateRequest.Description
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.Port != nil {
			mockState.Port = *updateRequest.Port
		}

		// Return updated state
		*stunServer = *mockState
	}).Maybe()

	// Mock the DeleteSTUNServer API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/stun_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityStunServer(t, client)
}

func testInfinityStunServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_stun_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "name", "tf-test-stun-server-full"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "description", "tf-test STUN server description"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "address", "stun-full.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "port", "5349"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_stun_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "name", "tf-test-stun-server"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "address", "stun.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "port", "3478"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_stun_server_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_stun_server_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "name", "tf-test-stun-server"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "address", "stun.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "port", "3478"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_stun_server_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.tf-test-stun-server", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "name", "tf-test-stun-server-full"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "description", "tf-test STUN server description"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "address", "stun-full.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.tf-test-stun-server", "port", "5349"),
				),
			},
		},
	})
}
