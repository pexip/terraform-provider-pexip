/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityMjxEndpoint(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateMjxendpoint API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_endpoint/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_endpoint/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.MjxEndpoint{
		ID:                             123,
		ResourceURI:                    "/api/admin/configuration/v1/mjx_endpoint/123/",
		Name:                           "mjx_endpoint-test",
		Description:                    "Test MjxEndpoint",
		EndpointType:                   "polycom",
		RoomResourceEmail:              "test@example.com",
		MjxEndpointGroup:               test.StringPtr("test-value"),
		APIAddress:                     test.StringPtr("test-server.example.com"),
		APIUsername:                    test.StringPtr("mjx_endpoint-test"),
		APIPassword:                    test.StringPtr("test-value"),
		UseHTTPS:                       "yes",
		VerifyCert:                     "yes",
		PolyUsername:                   test.StringPtr("mjx_endpoint-test"),
		PolyPassword:                   test.StringPtr("test-value"),
		PolyRaiseAlarmsForThisEndpoint: true,
		WebexDeviceID:                  test.StringPtr("test-value"),
	}

	// Mock the GetMjxendpoint API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mjx_endpoint := args.Get(2).(*config.MjxEndpoint)
		*mjx_endpoint = *mockState
	}).Maybe()

	// Mock the UpdateMjxendpoint API call
	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MjxEndpointUpdateRequest)
		mjx_endpoint := args.Get(3).(*config.MjxEndpoint)

		// Update mock state based on request
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.EndpointType != "" {
			mockState.EndpointType = updateRequest.EndpointType
		}
		if updateRequest.RoomResourceEmail != "" {
			mockState.RoomResourceEmail = updateRequest.RoomResourceEmail
		}
		if updateRequest.MjxEndpointGroup != nil {
			mockState.MjxEndpointGroup = updateRequest.MjxEndpointGroup
		}
		if updateRequest.APIAddress != nil {
			mockState.APIAddress = updateRequest.APIAddress
		}
		if updateRequest.APIPassword != nil {
			mockState.APIPassword = updateRequest.APIPassword
		}
		if updateRequest.UseHTTPS != "" {
			mockState.UseHTTPS = updateRequest.UseHTTPS
		}
		if updateRequest.VerifyCert != "" {
			mockState.VerifyCert = updateRequest.VerifyCert
		}
		if updateRequest.PolyPassword != nil {
			mockState.PolyPassword = updateRequest.PolyPassword
		}
		if updateRequest.PolyRaiseAlarmsForThisEndpoint != nil {
			mockState.PolyRaiseAlarmsForThisEndpoint = *updateRequest.PolyRaiseAlarmsForThisEndpoint
		}
		if updateRequest.WebexDeviceID != nil {
			mockState.WebexDeviceID = updateRequest.WebexDeviceID
		}

		// Return updated state
		*mjx_endpoint = *mockState
	}).Maybe()

	// Mock the DeleteMjxendpoint API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/mjx_endpoint/123/"
	}), mock.Anything).Return(nil)

	testInfinityMjxEndpoint(t, client)
}

func testInfinityMjxEndpoint(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "name", "mjx_endpoint-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "description", "Test MjxEndpoint"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "api_username", "mjx_endpoint-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "poly_username", "mjx_endpoint-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "poly_raise_alarms_for_this_endpoint", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "name", "mjx_endpoint-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "description", "Updated Test MjxEndpoint"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "api_username", "mjx_endpoint-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "poly_username", "mjx_endpoint-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.mjx_endpoint-test", "poly_raise_alarms_for_this_endpoint", "false"),
				),
			},
		},
	})
}
