/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityRoleMapping(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateRolemapping API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/role_mapping/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/role_mapping/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.RoleMapping{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/role_mapping/123/",
		Name:        "role_mapping-test",
		Source:      "saml_attribute",
		Value:       "test-value",
	}

	// Mock the GetRolemapping API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/role_mapping/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		role_mapping := args.Get(2).(*config.RoleMapping)
		*role_mapping = *mockState
	}).Maybe()

	// Mock the UpdateRolemapping API call
	client.On("PutJSON", mock.Anything, "configuration/v1/role_mapping/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.RoleMappingUpdateRequest)
		role_mapping := args.Get(3).(*config.RoleMapping)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		mockState.Source = updateRequest.Source
		mockState.Value = updateRequest.Value

		// Return updated state
		*role_mapping = *mockState
	}).Maybe()

	// Mock the DeleteRolemapping API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/role_mapping/123/"
	}), mock.Anything).Return(nil)

	testInfinityRoleMapping(t, client)
}

func testInfinityRoleMapping(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_mapping_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role_mapping.role_mapping-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role_mapping.role_mapping-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role_mapping.role_mapping-test", "name", "role_mapping-test"),
					resource.TestCheckResourceAttr("pexip_infinity_role_mapping.role_mapping-test", "source", "saml_attribute"),
					resource.TestCheckResourceAttr("pexip_infinity_role_mapping.role_mapping-test", "value", "test-value"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_mapping_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role_mapping.role_mapping-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role_mapping.role_mapping-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role_mapping.role_mapping-test", "name", "role_mapping-test"),
					resource.TestCheckResourceAttr("pexip_infinity_role_mapping.role_mapping-test", "source", "ldap_attribute"),
					resource.TestCheckResourceAttr("pexip_infinity_role_mapping.role_mapping-test", "value", "updated-value"),
				),
			},
		},
	})
}
