/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"fmt"
	"os"
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
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateRole API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/role/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/role/", mock.Anything, mock.Anything).Return(createResponse, nil).Maybe()

	// Shared state for mocking - starts with full config
	mockState := &config.Role{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/role/123/",
		Name:        "tf-test-role",
		Permissions: []config.Permission{
			{ID: 1, Name: "permission1"},
			{ID: 2, Name: "permission2"},
		},
	}

	// Mock the GetRole API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/role/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		role := args.Get(3).(*config.Role)
		*role = *mockState
	}).Maybe()

	// Mock the UpdateRole API call
	client.On("PutJSON", mock.Anything, "configuration/v1/role/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.RoleUpdateRequest)
		role := args.Get(3).(*config.Role)

		// Update mock state based on request
		mockState.Name = updateRequest.Name

		// Handle permissions - if provided, update; if nil, clear
		if updateRequest.Permissions != nil {
			if len(updateRequest.Permissions) == 0 {
				mockState.Permissions = []config.Permission{}
			} else {
				converted := make([]config.Permission, len(updateRequest.Permissions))
				for i, perm := range updateRequest.Permissions {
					// Extract ID from permission URI (e.g., "/api/admin/configuration/v1/permission/1/" -> 1)
					var id int
					fmt.Sscanf(perm, "/api/admin/configuration/v1/permission/%d/", &id)
					converted[i] = config.Permission{ID: id, Name: fmt.Sprintf("permission%d", id)}
				}
				mockState.Permissions = converted
			}
		}

		// Return updated state
		*role = *mockState
	}).Maybe()

	// Mock the DeleteRole API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/role/123/"
	}), mock.Anything).Return(nil).Maybe()

	testInfinityRole(t, client)
}

func testInfinityRole(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "name", "tf-test-role"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "permissions.#", "2"),
				),
			},
			// Step 2: Update to min config (clearing permissions)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "name", "tf-test-role"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "permissions.#", "0"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_role_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "name", "tf-test-role"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "permissions.#", "0"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_role_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_role.tf-test-role", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "name", "tf-test-role"),
					resource.TestCheckResourceAttr("pexip_infinity_role.tf-test-role", "permissions.#", "2"),
				),
			},
		},
	})
}
