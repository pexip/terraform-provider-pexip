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
		Name:        "identity_provider_attribute-test",
		Description: "Test IdentityProviderAttribute",
	}

	// Mock the GetIdentityproviderattribute API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/identity_provider_attribute/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		identity_provider_attribute := args.Get(2).(*config.IdentityProviderAttribute)
		*identity_provider_attribute = *mockState
	}).Maybe()

	// Mock the UpdateIdentityproviderattribute API call
	client.On("PutJSON", mock.Anything, "configuration/v1/identity_provider_attribute/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.IdentityProviderAttributeUpdateRequest)
		identity_provider_attribute := args.Get(3).(*config.IdentityProviderAttribute)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}

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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "name", "identity_provider_attribute-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "description", "Test IdentityProviderAttribute"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_identity_provider_attribute_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "name", "identity_provider_attribute-test"),
					resource.TestCheckResourceAttr("pexip_infinity_identity_provider_attribute.identity_provider_attribute-test", "description", "Updated Test IdentityProviderAttribute"),
				),
			},
		},
	})
}
