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

func TestInfinityGMSAccessToken(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateGmsaccesstoken API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/gms_access_token/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/gms_access_token/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.GMSAccessToken{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/gms_access_token/123/",
		Name:        "gms_access_token-test",
		Token:       "test-value",
	}

	// Mock the GetGmsaccesstoken API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/gms_access_token/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		gms_access_token := args.Get(3).(*config.GMSAccessToken)
		*gms_access_token = *mockState
	}).Maybe()

	// Mock the UpdateGmsaccesstoken API call
	client.On("PutJSON", mock.Anything, "configuration/v1/gms_access_token/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.GMSAccessTokenUpdateRequest)
		gms_access_token := args.Get(3).(*config.GMSAccessToken)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Token != "" {
			mockState.Token = updateRequest.Token
		}

		// Return updated state
		*gms_access_token = *mockState
	}).Maybe()

	// Mock the DeleteGmsaccesstoken API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/gms_access_token/123/"
	}), mock.Anything).Return(nil)

	testInfinityGMSAccessToken(t, client)
}

func testInfinityGMSAccessToken(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gms_access_token_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_access_token.gms_access_token-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_access_token.gms_access_token-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gms_access_token.gms_access_token-test", "name", "gms_access_token-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_gms_access_token_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_access_token.gms_access_token-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_gms_access_token.gms_access_token-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_gms_access_token.gms_access_token-test", "name", "gms_access_token-test"),
				),
			},
		},
	})
}
