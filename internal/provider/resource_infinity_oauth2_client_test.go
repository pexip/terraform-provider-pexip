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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityOAuth2Client(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateOauth2Client API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/oauth2_client/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/oauth2_client/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.OAuth2Client{
		ResourceURI: "/api/admin/configuration/v1/oauth2_client/123/",
		ClientID:    "123",
		ClientName:  "oauth2_client-test",
		Role:        "test-value",
	}

	// Mock the GetOauth2Client API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/oauth2_client/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		oauth2_client := args.Get(3).(*config.OAuth2Client)
		*oauth2_client = *mockState
	}).Maybe()

	// Mock the UpdateOauth2Client API call
	client.On("PutJSON", mock.Anything, "configuration/v1/oauth2_client/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.OAuth2ClientUpdateRequest)
		oauth2_client := args.Get(3).(*config.OAuth2Client)

		// Update mock state based on request
		mockState.ClientName = updateRequest.ClientName
		mockState.Role = updateRequest.Role

		// Return updated state
		*oauth2_client = *mockState
	}).Maybe()

	// Mock the DeleteOauth2Client API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/oauth2_client/123/", mock.Anything).Return(nil)

	testInfinityOAuth2Client(t, client)
}

func testInfinityOAuth2Client(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_oauth2_client_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "client_id"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "client_name", "oauth2_client-test"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "role", "test-value"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_oauth2_client_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "client_id"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "client_name", "oauth2_client-test"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "role", "updated-value"),
				),
			},
		},
	})
}
