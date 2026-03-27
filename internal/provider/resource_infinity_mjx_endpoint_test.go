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

	client := infinity.NewClientMock()

	mockState := &config.MjxEndpoint{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_endpoint/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_endpoint/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxEndpointCreateRequest)
		*mockState = config.MjxEndpoint{
			ID:                             123,
			ResourceURI:                    "/api/admin/configuration/v1/mjx_endpoint/123/",
			Name:                           createReq.Name,
			Description:                    createReq.Description,
			EndpointType:                   createReq.EndpointType,
			RoomResourceEmail:              createReq.RoomResourceEmail,
			MjxEndpointGroup:               createReq.MjxEndpointGroup,
			APIAddress:                     createReq.APIAddress,
			APIPort:                        createReq.APIPort,
			APIUsername:                    createReq.APIUsername,
			UseHTTPS:                       createReq.UseHTTPS,
			VerifyCert:                     createReq.VerifyCert,
			PolyUsername:                   createReq.PolyUsername,
			PolyRaiseAlarmsForThisEndpoint: createReq.PolyRaiseAlarmsForThisEndpoint,
			WebexDeviceID:                  createReq.WebexDeviceID,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mjxEndpoint := args.Get(3).(*config.MjxEndpoint)
		*mjxEndpoint = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxEndpointUpdateRequest)

		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.EndpointType = updateReq.EndpointType
		mockState.RoomResourceEmail = updateReq.RoomResourceEmail
		mockState.MjxEndpointGroup = updateReq.MjxEndpointGroup
		mockState.APIAddress = updateReq.APIAddress
		mockState.APIPort = updateReq.APIPort
		mockState.APIUsername = updateReq.APIUsername
		mockState.UseHTTPS = updateReq.UseHTTPS
		mockState.VerifyCert = updateReq.VerifyCert
		mockState.PolyUsername = updateReq.PolyUsername
		mockState.PolyRaiseAlarmsForThisEndpoint = updateReq.PolyRaiseAlarmsForThisEndpoint
		mockState.WebexDeviceID = updateReq.WebexDeviceID
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_endpoint/123/", mock.Anything).Return(nil)

	testInfinityMjxEndpoint(t, client)
}

func testInfinityMjxEndpoint(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "name", "tf-test mjx-endpoint full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "description", "Test MJX endpoint description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "endpoint_type", "CISCO"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "room_resource_email", "room2@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_address", "192.168.1.101"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_port", "443"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_username", "apiuser"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_password", "apipassword"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "use_https", "YES"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "verify_cert", "YES"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_username", "polyuser"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_password", "polypassword"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_raise_alarms_for_this_endpoint", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "webex_device_id", "device-123"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "name", "tf-test mjx-endpoint min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "endpoint_type", "CISCO"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "room_resource_email", "room@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "use_https", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "verify_cert", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_raise_alarms_for_this_endpoint", "true"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "api_port"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "api_username"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_username"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "webex_device_id"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "name", "tf-test mjx-endpoint min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "endpoint_type", "CISCO"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "room_resource_email", "room@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "use_https", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "verify_cert", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_raise_alarms_for_this_endpoint", "true"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "api_port"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "api_username"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_username"),
					resource.TestCheckNoResourceAttr("pexip_infinity_mjx_endpoint.test", "webex_device_id"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_endpoint_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_endpoint.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "name", "tf-test mjx-endpoint full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "description", "Test MJX endpoint description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "endpoint_type", "CISCO"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "room_resource_email", "room2@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_address", "192.168.1.101"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_port", "443"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_username", "apiuser"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "api_password", "apipassword"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "use_https", "YES"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "verify_cert", "YES"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_username", "polyuser"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_password", "polypassword"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "poly_raise_alarms_for_this_endpoint", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_endpoint.test", "webex_device_id", "device-123"),
				),
			},
		},
	})
}
