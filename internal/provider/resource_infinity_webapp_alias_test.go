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

func TestInfinityWebappAlias(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateWebappalias API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_alias/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/webapp_alias/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.WebappAlias{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/webapp_alias/123/",
		Slug:        "test-alias",
		Description: "Test WebappAlias",
		WebappType:  "webapp1",
		IsEnabled:   true,
	}

	// Mock the GetWebappalias API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/webapp_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		webapp_alias := args.Get(3).(*config.WebappAlias)
		*webapp_alias = *mockState
	}).Maybe()

	// Mock the UpdateWebappalias API call
	client.On("PutJSON", mock.Anything, "configuration/v1/webapp_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.WebappAliasUpdateRequest)
		webapp_alias := args.Get(3).(*config.WebappAlias)

		// Update mock state based on request
		if updateReq.Slug != "" {
			mockState.Slug = updateReq.Slug
		}
		if updateReq.Description != "" {
			mockState.Description = updateReq.Description
		}
		if updateReq.WebappType != "" {
			mockState.WebappType = updateReq.WebappType
		}
		if updateReq.IsEnabled != nil {
			mockState.IsEnabled = *updateReq.IsEnabled
		}
		if updateReq.Bundle != nil {
			mockState.Bundle = updateReq.Bundle
		}
		if updateReq.Branding != nil {
			mockState.Branding = updateReq.Branding
		}

		// Return updated state
		*webapp_alias = *mockState
	}).Maybe()

	// Mock the DeleteWebappalias API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/webapp_alias/123/"
	}), mock.Anything).Return(nil)

	testInfinityWebappAlias(t, client)
}

func testInfinityWebappAlias(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_alias_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.webapp_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.webapp_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.webapp_alias-test", "description", "Test WebappAlias"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.webapp_alias-test", "is_enabled", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_alias_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.webapp_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.webapp_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.webapp_alias-test", "slug", "updated-alias"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.webapp_alias-test", "description", "Updated Test WebappAlias"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.webapp_alias-test", "webapp_type", "webapp2"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.webapp_alias-test", "is_enabled", "false"),
				),
			},
		},
	})
}
