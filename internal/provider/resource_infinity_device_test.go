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

func TestInfinityDevice(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock identity provider group creation
	idpGroupCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider_group/456/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider_group/", mock.Anything, mock.Anything).Return(idpGroupCreateResponse, nil)

	// Mock identity provider group state
	idpGroupState := &config.IdentityProviderGroup{
		ID:          456,
		ResourceURI: "/api/admin/configuration/v1/identity_provider_group/456/",
		Name:        "tf-test-identity-provider-group",
		Description: "Test Identity Provider Group for Device",
	}

	// Mock GetIdentityProviderGroup for read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider_group/456/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		group := args.Get(3).(*config.IdentityProviderGroup)
		*group = *idpGroupState
	}).Maybe()

	// Mock UpdateIdentityProviderGroup
	client.On("PutJSON", mock.Anything, "configuration/v1/identity_provider_group/456/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.IdentityProviderGroupUpdateRequest)
		group := args.Get(3).(*config.IdentityProviderGroup)
		idpGroupState.Name = updateRequest.Name
		idpGroupState.Description = updateRequest.Description
		*group = *idpGroupState
	}).Maybe()

	// Mock DeleteIdentityProviderGroup
	client.On("DeleteJSON", mock.Anything, "configuration/v1/identity_provider_group/456/", mock.Anything).Return(nil)

	// Mock the CreateDevice API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/device/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/device/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking - note: password is not returned by API
	idpGroupURI := "/api/admin/configuration/v1/identity_provider_group/456/"
	mockState := &config.Device{
		ID:                          123,
		ResourceURI:                 "/api/admin/configuration/v1/device/123/",
		Alias:                       "tf-test-device",
		Description:                 "Test Device Description",
		Username:                    "tf-test-user",
		Password:                    "", // API doesn't return password
		PrimaryOwnerEmailAddress:    "tf-test@example.com",
		EnableSIP:                   true,
		EnableH323:                  true,
		EnableInfinityConnectNonSSO: true,
		EnableInfinityConnectSSO:    true,
		EnableStandardSSO:           true,
		SSOIdentityProviderGroup:    &idpGroupURI,
		Tag:                         "tf-test-tag",
		SyncTag:                     "tf-test-sync-tag",
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

		// Update mock state - always update fields from request (except password, which API doesn't return)
		mockState.Alias = updateRequest.Alias
		mockState.Description = updateRequest.Description
		mockState.Username = updateRequest.Username
		// Note: password is updated internally but not returned by API
		mockState.PrimaryOwnerEmailAddress = updateRequest.PrimaryOwnerEmailAddress
		mockState.Tag = updateRequest.Tag
		mockState.SyncTag = updateRequest.SyncTag

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
		if updateRequest.SSOIdentityProviderGroup != nil {
			mockState.SSOIdentityProviderGroup = updateRequest.SSOIdentityProviderGroup
		}

		// Return updated state (password remains empty as API doesn't return it)
		*device = *mockState
	}).Maybe()

	// Mock the DeleteDevice API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/device/123/"
	}), mock.Anything).Return(nil)

	testInfinityDevice(t, client)
}

func testInfinityDevice(t *testing.T, client InfinityClient) {
	// Test 1 & 2: Create with full config, update to min config, then delete
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Test 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_device_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "alias", "tf-test-device"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "description", "Test Device Description"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "password", "tf-test-pass"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "primary_owner_email_address", "tf-test@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_sip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_h323", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_infinity_connect_non_sso", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_infinity_connect_sso", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_standard_sso", "true"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "sso_identity_provider_group"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "tag", "tf-test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "sync_tag", "tf-test-sync-tag"),
				),
			},
			// Test 2: Update to min config (then delete)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_device_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "alias", "tf-test-device"),
				),
			},
		},
	})

	// Test 3 & 4: Create with min config, update to full config
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Test 3: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_device_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "alias", "tf-test-device"),
				),
			},
			// Test 4: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_device_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "alias", "tf-test-device"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "description", "Test Device Description"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "username", "tf-test-user"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "password", "tf-test-pass"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "primary_owner_email_address", "tf-test@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_sip", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_h323", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_infinity_connect_non_sso", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_infinity_connect_sso", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "enable_standard_sso", "true"),
					resource.TestCheckResourceAttrSet("pexip_infinity_device.tf-test-device", "sso_identity_provider_group"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "tag", "tf-test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_device.tf-test-device", "sync_tag", "tf-test-sync-tag"),
				),
			},
		},
	})
}
