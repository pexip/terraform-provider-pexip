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

func TestInfinityRole(t *testing.T) {
	t.Parallel()

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateRole API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/role/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/role/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.Role{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/role/123/",
		Name:        "role-test",
		Permissions: []string{}, // Empty permissions list
	}

	// Mock the GetRole API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/role/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		role := args.Get(2).(*config.Role)
		*role = *mockState
	}).Maybe()

	// Mock the UpdateRole API call
	client.On("PutJSON", mock.Anything, "configuration/v1/role/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.RoleUpdateRequest)
		role := args.Get(3).(*config.Role)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		if updateRequest.Permissions != nil {
			mockState.Permissions = updateRequest.Permissions
		}

		// Return updated state
		*role = *mockState
	}).Maybe()

	// Mock the DeleteRole API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/role/123/"
	}), mock.Anything).Return(nil)

	testInfinityRole(t, client)
}

func testInfinityRole(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role.role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role.role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role.role-test", "name", "role-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role.role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role.role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role.role-test", "name", "role-test"),
				),
			},
		},
	})
}
