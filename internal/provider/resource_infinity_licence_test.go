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

func TestInfinityLicence(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateLicence API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/licence/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/licence/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.Licence{
		FulfillmentID: "test-fulfillment-123",
		EntitlementID: "test-value",
		OfflineMode:   true,
		ResourceURI:   "/api/admin/configuration/v1/licence/123/",
	}

	// Mock the ListLicences API call (needed after creation to find fulfillment ID)
	listResponse := &config.LicenceListResponse{
		Objects: []config.Licence{*mockState},
	}
	client.On("GetJSON", mock.Anything, "configuration/v1/licence/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(2).(*config.LicenceListResponse)
		*response = *listResponse
	})

	// Mock the GetLicence API call for Read operations (both paths needed)
	client.On("GetJSON", mock.Anything, "configuration/v1/licence/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		license := args.Get(2).(*config.Licence)
		*license = *mockState
	}).Maybe()
	client.On("GetJSON", mock.Anything, "configuration/v1/licence/test-fulfillment-123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		license := args.Get(2).(*config.Licence)
		*license = *mockState
	}).Maybe()

	// Licence doesn't support update operations

	// Mock the DeleteLicence API call (uses FulfillmentID from resource)
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/licence/test-fulfillment-123/"
	}), mock.Anything).Return(nil)

	testInfinityLicence(t, client)
}

func testInfinityLicence(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_licence_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_licence.licence-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_licence.licence-test", "entitlement_id", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_licence.licence-test", "offline_mode", "true"),
				),
			},
			// Licence doesn't support updates, so only test creation/read
		},
	})
}
