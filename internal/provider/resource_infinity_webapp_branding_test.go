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

func TestInfinityWebappBranding(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateWebappbranding API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_branding/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/webapp_branding/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.WebappBranding{
		ResourceURI:  "/api/admin/configuration/v1/webapp_branding/123/",
		Name:         "webapp_branding-test",
		Description:  "Test WebappBranding",
		UUID:         "test-value",
		WebappType:   "pexapp",
		IsDefault:    true,
		BrandingFile: "test-value",
	}

	// Mock the GetWebappbranding API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/webapp_branding/webapp_branding-test/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_branding := args.Get(3).(*config.WebappBranding)
		*webapp_branding = *mockState
	}).Maybe()

	// Mock the UpdateWebappbranding API call
	client.On("PutJSON", mock.Anything, "configuration/v1/webapp_branding/webapp_branding-test/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.WebappBrandingUpdateRequest)
		webapp_branding := args.Get(3).(*config.WebappBranding)

		// Update mock state based on request
		if updateReq.Description != "" {
			mockState.Description = updateReq.Description
		}
		if updateReq.UUID != "" {
			mockState.UUID = updateReq.UUID
		}
		if updateReq.WebappType != "" {
			mockState.WebappType = updateReq.WebappType
		}
		if updateReq.IsDefault != nil {
			mockState.IsDefault = *updateReq.IsDefault
		}
		if updateReq.BrandingFile != "" {
			mockState.BrandingFile = updateReq.BrandingFile
		}

		// Return updated state
		*webapp_branding = *mockState
	}).Maybe()

	// Mock the DeleteWebappbranding API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/webapp_branding/webapp_branding-test/"
	}), mock.Anything).Return(nil)

	testInfinityWebappBranding(t, client)
}

func testInfinityWebappBranding(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "name", "webapp_branding-test"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "description", "Test WebappBranding"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "is_default", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_branding_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_branding.webapp_branding-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "name", "webapp_branding-test"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "description", "Updated Test WebappBranding"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "uuid", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "webapp_type", "management"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "is_default", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_branding.webapp_branding-test", "branding_file", "updated-value"),
				),
			},
		},
	})
}
