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

func TestInfinityWebappAlias(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.WebappAlias{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/webapp_alias/123/",
		Slug:        "tf-test-alias-full",
		Description: "tf-test webapp alias description",
		WebappType:  "webapp3",
		IsEnabled:   true,
	}

	// Mock the CreateWebappalias API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/webapp_alias/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/webapp_alias/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.WebappAliasCreateRequest)
		// Update mock state based on create request
		mockState.Slug = createReq.Slug
		mockState.Description = createReq.Description
		mockState.WebappType = createReq.WebappType
		mockState.IsEnabled = createReq.IsEnabled
	}).Maybe()

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
		// Description can be empty string, so always update
		mockState.Description = updateReq.Description
		if updateReq.WebappType != "" {
			mockState.WebappType = updateReq.WebappType
		}
		if updateReq.IsEnabled != nil {
			mockState.IsEnabled = *updateReq.IsEnabled
		}
		// Note: Bundle and Branding are sent as string URIs in requests,
		// but returned as objects in responses. We don't update them in the mock
		// since the test doesn't use these fields.

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
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_alias_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "slug", "tf-test-alias-full"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "description", "tf-test webapp alias description"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "webapp_type", "webapp3"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "is_enabled", "true"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_alias_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "slug", "tf-test-alias"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "webapp_type", "webapp1"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "is_enabled", "false"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_webapp_alias_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_alias_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "slug", "tf-test-alias"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "webapp_type", "webapp1"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "is_enabled", "false"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_webapp_alias_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_webapp_alias.tf-test-webapp-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "slug", "tf-test-alias-full"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "description", "tf-test webapp alias description"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "webapp_type", "webapp3"),
					resource.TestCheckResourceAttr("pexip_infinity_webapp_alias.tf-test-webapp-alias", "is_enabled", "true"),
				),
			},
		},
	})
}
