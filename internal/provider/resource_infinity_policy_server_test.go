/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityPolicyServer(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreatePolicyServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/policy_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/policy_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.PolicyServer{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/policy_server/123/",
		Name:        "test-policy-server",
		Description: "Test Policy Server",
		URL:         "https://test-policy-server.dev.pexip.network",
		Username:    "testuser",
		Password:    "testpassword",
	}

	// Mock the GetPolicyServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/policy_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		policyServer := args.Get(2).(*config.PolicyServer)
		*policyServer = *mockState
	}).Maybe()

	// Mock the UpdatePolicyServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/policy_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.PolicyServerUpdateRequest)
		policyServer := args.Get(3).(*config.PolicyServer)

		// Update mock state
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.URL != "" {
			mockState.URL = updateRequest.URL
		}
		if updateRequest.Username != "" {
			mockState.Username = updateRequest.Username
		}
		if updateRequest.Password != "" {
			mockState.Password = updateRequest.Password
		}

		// Return updated state
		*policyServer = *mockState
	}).Maybe()

	// Mock the DeletePolicyServer API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/policy_server/123/"
	}), mock.Anything).Return(nil)

	testInfinityPolicyServer(t, client)
}

func testInfinityPolicyServer(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_policy_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.policy-server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.policy-server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "name", "test-policy-server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "description", "Test Policy Server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "url", "https://test-policy-server.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "username", "testuser"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "password", "testpassword"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_policy_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.policy-server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_policy_server.policy-server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "name", "test-policy-server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "description", "Test Policy Server"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "url", "https://test-policy-server.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "username", "testuser"),
					resource.TestCheckResourceAttr("pexip_infinity_policy_server.policy-server-test", "password", "updatedpassword"),
				),
			},
		},
	})
}
