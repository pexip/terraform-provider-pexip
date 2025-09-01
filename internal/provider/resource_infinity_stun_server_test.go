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

func TestInfinityStunServer(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSTUNServer API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/stun_server/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/stun_server/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.STUNServer{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/stun_server/123/",
		Name:        "stun-server-test",
		Description: "Test STUN server",
		Address:     "test-stun-server.dev.pexip.network",
		Port:        8080,
	}

	// Mock the GetSTUNServer API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/stun_server/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		stunServer := args.Get(2).(*config.STUNServer)
		*stunServer = *mockState
	}).Maybe()

	// Mock the UpdateSTUNServer API call
	client.On("PutJSON", mock.Anything, "configuration/v1/stun_server/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.STUNServerUpdateRequest)
		stunServer := args.Get(3).(*config.STUNServer)

		// Update mock state
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
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
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_stun_server_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.stun-server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.stun-server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "name", "stun-server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "description", "Test STUN server"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "address", "test-stun-server.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "port", "8080"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_stun_server_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.stun-server-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_stun_server.stun-server-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "name", "stun-server-test"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "description", "Test STUN server"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "address", "test-stun-server.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_stun_server.stun-server-test", "port", "8081"),
				),
			},
		},
	})
}
