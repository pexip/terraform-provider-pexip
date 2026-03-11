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

func TestInfinityIdentityProviderAttribute(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateIdentityproviderattribute API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/identity_provider_attribute/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/identity_provider_attribute/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.IdentityProviderAttribute{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/identity_provider_attribute/123/",
		Name:        "tf-test-identity-provider-attribute",
		Description: "Test Identity Provider Attribute Description",
	}

	// Mock the GetIdentityproviderattribute API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider_attribute/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		identity_provider_attribute := args.Get(3).(*config.IdentityProviderAttribute)
		*identity_provider_attribute = *mockState
	}).Maybe()

	// Mock the UpdateIdentityproviderattribute API call
	client.On("PutJSON", mock.Anything, "configuration/v1/identity_provider_attribute/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.IdentityProviderAttributeUpdateRequest)
		identity_provider_attribute := args.Get(3).(*config.IdentityProviderAttribute)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description

		// Return updated state
		*identity_provider_attribute = *mockState
	}).Maybe()

	// Mock the DeleteIdentityproviderattribute API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/identity_provider_attribute/123/"
	}), mock.Anything).Return(nil)

	testInfinityIdentityProviderAttribute(t, client)
}

func testInfinityIdentityProviderAttribute(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "name", "tf-test-identity-provider-attribute"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "description", "Test Identity Provider Attribute Description"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "name", "tf-test-identity-provider-attribute"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "name", "tf-test-identity-provider-attribute"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "name", "tf-test-identity-provider-attribute"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.tf-test-identity-provider-attribute", "description", "Test Identity Provider Attribute Description"),
				),
			},
		},
	})
}
