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

func TestInfinityLocation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateLocation API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/location/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/location/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.Location{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/location/123/",
		Name:        "location-test",
		Description: "Test Location",
	}

	// Mock the GetLocation API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/location/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		location := args.Get(3).(*config.Location)
		*location = *mockState
	}).Maybe()

	// Mock the UpdateLocation API call
	client.On("PutJSON", mock.Anything, "configuration/v1/location/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.LocationUpdateRequest)
		location := args.Get(3).(*config.Location)

		// Update mock state
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}

		// Return updated state
		*location = *mockState
	}).Maybe()

	// Mock the DeleteLocation API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/location/123/"
	}), mock.Anything).Return(nil)

	testInfinityLocation(t, client)
}

func testInfinityLocation(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_location_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_location.location-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_location.location-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_location.location-test", "name", "location-test"),
					resource.TestCheckResourceAttr("pexip_infinity_location.location-test", "description", "Test Location"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_location_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_location.location-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_location.location-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_location.location-test", "name", "location-test"),
					resource.TestCheckResourceAttr("pexip_infinity_location.location-test", "description", "Updated Test Location"),
				),
			},
		},
	})
}
