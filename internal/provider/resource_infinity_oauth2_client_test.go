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

func TestInfinityOAuth2Client(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateOauth2Client API call
	createResponseBody := `{
		"resource_uri": "/api/admin/configuration/v1/oauth2_client/123/",
		"client_id": "test-oauth2-client-id",
		"client_name": "tf-test oauth2_client RW",
		"role": "/api/admin/configuration/v1/role/1/",
		"private_key_jwt": "test-private-key-jwt"
	}`
	createResponse := &types.PostResponse{
		Body:        []byte(createResponseBody),
		ResourceURI: "/api/admin/configuration/v1/oauth2_client/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/oauth2_client/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.OAuth2Client{
		ResourceURI:   "/api/admin/configuration/v1/oauth2_client/123/",
		ClientID:      "123",
		ClientName:    "tf-test oauth2_client RW",
		Role:          "/api/admin/configuration/v1/role/1/",
		PrivateKeyJWT: "test-private-key-jwt",
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
		if updateRequest.ClientName != "" {
			mockState.ClientName = updateRequest.ClientName
		}
		if updateRequest.Role != "" {
			mockState.Role = updateRequest.Role
		}

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
				Config: test.LoadTestFolder(t, "resource_infinity_oauth2_client_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "client_id"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "client_name", "tf-test oauth2_client RW"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "role", "/api/admin/configuration/v1/role/1/"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "private_key_jwt"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_oauth2_client_full_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "client_id"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "client_name", "tf-test oauth2_client RO"),
					resource.TestCheckResourceAttr("pexip_infinity_oauth2_client.oauth2_client-test", "role", "/api/admin/configuration/v1/role/2/"),
					resource.TestCheckResourceAttrSet("pexip_infinity_oauth2_client.oauth2_client-test", "private_key_jwt"),
				),
			},
		},
	})
}
