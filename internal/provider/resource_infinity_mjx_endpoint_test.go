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

func TestInfinityMjxEndpoint(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for cisco endpoint
	mockStateCisco := &config.MjxEndpoint{}

	// Shared state for poly endpoint
	mockStatePoly := &config.MjxEndpoint{}

	// Mock the CreateMjxendpoint API call for cisco endpoint
	ciscoCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_endpoint/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_endpoint/", mock.MatchedBy(func(req *config.MjxEndpointCreateRequest) bool {
		return req.Name == "tf-test Cisco mjx-endpoint min" || req.Name == "tf-test Cisco mjx-endpoint full"
	}), mock.Anything).Return(ciscoCreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxEndpointCreateRequest)
		*mockStateCisco = config.MjxEndpoint{
			ID:                             123,
			ResourceURI:                    "/api/admin/configuration/v1/mjx_endpoint/123/",
			Name:                           createReq.Name,
			Description:                    createReq.Description,
			EndpointType:                   createReq.EndpointType,
			RoomResourceEmail:              createReq.RoomResourceEmail,
			MjxEndpointGroup:               createReq.MjxEndpointGroup,
			APIAddress:                     createReq.APIAddress,
			APIUsername:                    createReq.APIUsername,
			APIPassword:                    createReq.APIPassword,
			UseHTTPS:                       createReq.UseHTTPS,
			VerifyCert:                     createReq.VerifyCert,
			PolyUsername:                   createReq.PolyUsername,
			PolyPassword:                   createReq.PolyPassword,
			PolyRaiseAlarmsForThisEndpoint: createReq.PolyRaiseAlarmsForThisEndpoint,
			WebexDeviceID:                  createReq.WebexDeviceID,
		}
	})

	// Mock the CreateMjxendpoint API call for poly endpoint
	polyCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_endpoint/124/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_endpoint/", mock.MatchedBy(func(req *config.MjxEndpointCreateRequest) bool {
		return req.Name == "tf-test Poly mjx-endpoint min" || req.Name == "tf-test Poly mjx-endpoint full"
	}), mock.Anything).Return(polyCreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxEndpointCreateRequest)
		*mockStatePoly = config.MjxEndpoint{
			ID:                             124,
			ResourceURI:                    "/api/admin/configuration/v1/mjx_endpoint/124/",
			Name:                           createReq.Name,
			Description:                    createReq.Description,
			EndpointType:                   createReq.EndpointType,
			RoomResourceEmail:              createReq.RoomResourceEmail,
			MjxEndpointGroup:               createReq.MjxEndpointGroup,
			APIAddress:                     createReq.APIAddress,
			APIUsername:                    createReq.APIUsername,
			APIPassword:                    createReq.APIPassword,
			UseHTTPS:                       createReq.UseHTTPS,
			VerifyCert:                     createReq.VerifyCert,
			PolyUsername:                   createReq.PolyUsername,
			PolyPassword:                   createReq.PolyPassword,
			PolyRaiseAlarmsForThisEndpoint: createReq.PolyRaiseAlarmsForThisEndpoint,
			WebexDeviceID:                  createReq.WebexDeviceID,
		}
	})

	// Mock the GetMjxendpoint API call for Read operations - cisco
	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mjx_endpoint := args.Get(3).(*config.MjxEndpoint)
		*mjx_endpoint = *mockStateCisco
	}).Maybe()

	// Mock the GetMjxendpoint API call for Read operations - poly
	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_endpoint/124/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mjx_endpoint := args.Get(3).(*config.MjxEndpoint)
		*mjx_endpoint = *mockStatePoly
	}).Maybe()

	// Mock the UpdateMjxendpoint API call - cisco
	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MjxEndpointUpdateRequest)
		mjx_endpoint := args.Get(3).(*config.MjxEndpoint)

		// Update mock state based on request
		mockStateCisco.Name = updateRequest.Name
		mockStateCisco.Description = updateRequest.Description
		mockStateCisco.EndpointType = updateRequest.EndpointType
		mockStateCisco.RoomResourceEmail = updateRequest.RoomResourceEmail
		mockStateCisco.MjxEndpointGroup = updateRequest.MjxEndpointGroup
		mockStateCisco.APIAddress = updateRequest.APIAddress
		mockStateCisco.APIUsername = updateRequest.APIUsername
		mockStateCisco.APIPassword = updateRequest.APIPassword
		mockStateCisco.UseHTTPS = updateRequest.UseHTTPS
		mockStateCisco.VerifyCert = updateRequest.VerifyCert
		mockStateCisco.PolyUsername = updateRequest.PolyUsername
		mockStateCisco.PolyPassword = updateRequest.PolyPassword
		if updateRequest.PolyRaiseAlarmsForThisEndpoint != nil {
			mockStateCisco.PolyRaiseAlarmsForThisEndpoint = *updateRequest.PolyRaiseAlarmsForThisEndpoint
		}
		mockStateCisco.WebexDeviceID = updateRequest.WebexDeviceID

		// Return updated state
		*mjx_endpoint = *mockStateCisco
	}).Maybe()

	// Mock the UpdateMjxendpoint API call - poly
	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_endpoint/124/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MjxEndpointUpdateRequest)
		mjx_endpoint := args.Get(3).(*config.MjxEndpoint)

		// Update mock state based on request
		mockStatePoly.Name = updateRequest.Name
		mockStatePoly.Description = updateRequest.Description
		mockStatePoly.EndpointType = updateRequest.EndpointType
		mockStatePoly.RoomResourceEmail = updateRequest.RoomResourceEmail
		mockStatePoly.MjxEndpointGroup = updateRequest.MjxEndpointGroup
		mockStatePoly.APIAddress = updateRequest.APIAddress
		mockStatePoly.APIUsername = updateRequest.APIUsername
		mockStatePoly.APIPassword = updateRequest.APIPassword
		mockStatePoly.UseHTTPS = updateRequest.UseHTTPS
		mockStatePoly.VerifyCert = updateRequest.VerifyCert
		mockStatePoly.PolyUsername = updateRequest.PolyUsername
		mockStatePoly.PolyPassword = updateRequest.PolyPassword
		if updateRequest.PolyRaiseAlarmsForThisEndpoint != nil {
			mockStatePoly.PolyRaiseAlarmsForThisEndpoint = *updateRequest.PolyRaiseAlarmsForThisEndpoint
		}
		mockStatePoly.WebexDeviceID = updateRequest.WebexDeviceID

		// Return updated state
		*mjx_endpoint = *mockStatePoly
	}).Maybe()

	// Mock the DeleteMjxendpoint API call - cisco
	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything).Return(nil)

	// Mock the DeleteMjxendpoint API call - poly
	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_endpoint/124/", mock.Anything).Return(nil)

	testInfinityMjxEndpoint(t, client)
}

func testInfinityMjxEndpoint(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Test 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_full"),
				Check: resource.ComposeTestCheckFunc(
					// Cisco endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "room_resource_email"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "api_address"),
					// Poly endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "room_resource_email"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "api_address"),
				),
			},
			// Test 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_min"),
				Check: resource.ComposeTestCheckFunc(
					// Cisco endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "name"),
					// Poly endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "name"),
				),
			},
			// Test 3: Destroy and recreate with min config
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_min"),
				Destroy: true,
			},
			// Test 4: Recreate with min config (after destroy)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_min"),
				Check: resource.ComposeTestCheckFunc(
					// Cisco endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "name"),
					// Poly endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "name"),
				),
			},
			// Test 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_full"),
				Check: resource.ComposeTestCheckFunc(
					// Cisco endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "room_resource_email"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.cisco", "api_address"),
					// Poly endpoint
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "resource_id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "room_resource_email"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.poly", "api_address"),
				),
			},
		},
	})
}
