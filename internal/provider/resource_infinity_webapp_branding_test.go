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

func TestInfinityWebappBranding(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - UUID is generated and returned by the API
	createdUUID := "12345678-1234-1234-1234-123456789012"
	mockState := &config.WebappBranding{
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/" + createdUUID + "/",
		Name:        "webapp_branding-test",
		Description: "Test WebappBranding",
		UUID:        createdUUID,
		WebappType:  "webapp2",
		IsDefault:   false,
	}

	// Mock the CreateWebappBranding API call (uses PostMultipartFormWithFieldsAndResponseUUID)
	createResponse := &types.PostResponseWithUUID{
		Body:         []byte(""),
		ResourceUUID: "/api/admin/configuration/v1/webapp_branding/" + createdUUID + "/",
	}
	client.On("PostMultipartFormWithFieldsAndResponseUUID", mock.Anything, "configuration/v1/webapp_branding/", mock.Anything, "branding_file", mock.Anything, mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		fields := args.Get(2).(map[string]string)
		// Update mock state based on create request fields
		mockState.Name = fields["name"]
		mockState.Description = fields["description"]
		mockState.WebappType = fields["webapp_type"]
	})

	// Mock the GetWebappBranding API call for Read operations (uses UUID)
	client.On("GetJSON", mock.Anything, "configuration/v1/webapp_branding/"+createdUUID+"/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_branding := args.Get(3).(*config.WebappBranding)
		*webapp_branding = *mockState
	}).Maybe()

	// Mock the UpdateWebappBranding API call (uses PatchJSON, not multipart)
	client.On("PatchJSON", mock.Anything, "configuration/v1/webapp_branding/"+createdUUID+"/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.WebappBrandingUpdateRequest)
		result := args.Get(3).(*config.WebappBranding)

		// Update mock state based on update request
		if req.Description != "" {
			mockState.Description = req.Description
		}

		// Return updated state
		*result = *mockState
	}).Maybe()

	// Mock the DeleteWebappBranding API call (uses UUID)
	client.On("DeleteJSON", mock.Anything, "configuration/v1/webapp_branding/"+createdUUID+"/", mock.Anything).Return(nil)

	testInfinityWebappBranding(t, client)
}

func testInfinityWebappBranding(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "name", "webapp_branding-test"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "description", "Test WebappBranding"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "webapp_type", "webapp2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "branding_file"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "last_updated"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "uuid"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "name", "webapp_branding-test"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "description", "Updated Test WebappBranding"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "webapp_type", "webapp2"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "branding_file"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "last_updated"),
				),
			},
		},
	})
}
