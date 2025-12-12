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

func TestInfinityDevice(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateDevice API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/device/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/device/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.Device{
		ID:                          123,
		ResourceURI:                 "/api/admin/configuration/v1/device/123/",
		Alias:                       "device-test",
		Description:                 "Test Device",
		Username:                    "deviceuser",
		Password:                    "devicepass",
		PrimaryOwnerEmailAddress:    "owner@example.com",
		EnableSIP:                   true,
		EnableH323:                  false,
		EnableInfinityConnectNonSSO: true,
		EnableInfinityConnectSSO:    false,
		EnableStandardSSO:           false,
		Tag:                         "test-tag",
		SyncTag:                     "sync-tag",
	}

	// Mock the GetDevice API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/device/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		device := args.Get(3).(*config.Device)
		*device = *mockState
	}).Maybe()

	// Mock the UpdateDevice API call
	client.On("PutJSON", mock.Anything, "configuration/v1/device/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.DeviceUpdateRequest)
		device := args.Get(3).(*config.Device)

		// Update mock state
		mockState.Alias = updateRequest.Alias
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Username != "" {
			mockState.Username = updateRequest.Username
		}
		if updateRequest.Password != "" {
			mockState.Password = updateRequest.Password
		}
		if updateRequest.PrimaryOwnerEmailAddress != "" {
			mockState.PrimaryOwnerEmailAddress = updateRequest.PrimaryOwnerEmailAddress
		}
		if updateRequest.Tag != "" {
			mockState.Tag = updateRequest.Tag
		}
		if updateRequest.SyncTag != "" {
			mockState.SyncTag = updateRequest.SyncTag
		}
		if updateRequest.EnableSIP != nil {
			mockState.EnableSIP = *updateRequest.EnableSIP
		}
		if updateRequest.EnableH323 != nil {
			mockState.EnableH323 = *updateRequest.EnableH323
		}
		if updateRequest.EnableInfinityConnectNonSSO != nil {
			mockState.EnableInfinityConnectNonSSO = *updateRequest.EnableInfinityConnectNonSSO
		}
		if updateRequest.EnableInfinityConnectSSO != nil {
			mockState.EnableInfinityConnectSSO = *updateRequest.EnableInfinityConnectSSO
		}
		if updateRequest.EnableStandardSSO != nil {
			mockState.EnableStandardSSO = *updateRequest.EnableStandardSSO
		}

		// Return updated state
		*device = *mockState
	}).Maybe()

	// Mock the DeleteDevice API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/device/123/"
	}), mock.Anything).Return(nil)

	testInfinityDevice(t, client)
}

func testInfinityDevice(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_device_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_device.device-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.device-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "alias", "device-test"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "description", "Test Device"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "username", "deviceuser"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "password", "devicepass"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "primary_owner_email_address", "owner@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_sip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_h323", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_infinity_connect_non_sso", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_infinity_connect_sso", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_standard_sso", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "tag", "test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "sync_tag", "sync-tag"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_device_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_device.device-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.device-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "alias", "device-test"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "description", "Updated Test Device"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "username", "updateduser"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "password", "updatedpass"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "primary_owner_email_address", "updated@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_sip", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_h323", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_infinity_connect_non_sso", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_infinity_connect_sso", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "enable_standard_sso", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "tag", "updated-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_device.device-test", "sync_tag", "updated-sync-tag"),
				),
			},
		},
	})
}
