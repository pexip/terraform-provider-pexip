/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityLicence(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.Licence{
		FulfillmentID: "test-fulfillment-123",
		EntitlementID: "",
		ResourceURI:   "/api/admin/configuration/v1/licence/123/",
	}

	// Mock the CreateLicence API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/licence/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/licence/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.LicenceCreateRequest)
		// Store the entitlement ID from the create request (without spaces)
		mockState.EntitlementID = strings.ReplaceAll(req.EntitlementID, " ", "")
	})

	// Mock the ListLicences API call (needed after creation to find fulfillment ID)
	client.On("GetJSON", mock.Anything, "configuration/v1/licence/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		response := args.Get(3).(*config.LicenceListResponse)
		// Dynamically return the current state of mockState
		response.Objects = []config.Licence{*mockState}
	})

	// Mock the GetLicence API call for Read operations (both paths needed)
	client.On("GetJSON", mock.Anything, "configuration/v1/licence/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		license := args.Get(3).(*config.Licence)
		*license = *mockState
	}).Maybe()
	client.On("GetJSON", mock.Anything, "configuration/v1/licence/test-fulfillment-123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		license := args.Get(3).(*config.Licence)
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
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_licence"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_licence.licence-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_licence.licence-test", "entitlement_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_licence.licence-test", "offline_mode"),
				),
			},
			// Licence doesn't support updates, so only test creation/read
		},
	})
}
