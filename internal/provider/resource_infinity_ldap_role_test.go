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

func TestInfinityLdapRole(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateRole API call (needed because ldap role references roles)
	role1CreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/role/1/",
	}
	mockRole1 := &config.Role{}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/role/", mock.MatchedBy(func(req *config.RoleCreateRequest) bool {
		return req.Name == "tf-test role 1 for ldap role"
	}), mock.Anything).Return(role1CreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.RoleCreateRequest)
		*mockRole1 = config.Role{
			ID:          1,
			ResourceURI: "/api/admin/configuration/v1/role/1/",
			Name:        createReq.Name,
		}
	})

	role2CreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/role/2/",
	}
	mockRole2 := &config.Role{}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/role/", mock.MatchedBy(func(req *config.RoleCreateRequest) bool {
		return req.Name == "tf-test role 2 for ldap role"
	}), mock.Anything).Return(role2CreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.RoleCreateRequest)
		*mockRole2 = config.Role{
			ID:          2,
			ResourceURI: "/api/admin/configuration/v1/role/2/",
			Name:        createReq.Name,
		}
	})

	// Mock the GetRole API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/role/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		role := args.Get(3).(*config.Role)
		*role = *mockRole1
	}).Maybe()

	client.On("GetJSON", mock.Anything, "configuration/v1/role/2/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		role := args.Get(3).(*config.Role)
		*role = *mockRole2
	}).Maybe()

	// Mock the DeleteRole API call
	client.On("DeleteJSON", mock.Anything, "configuration/v1/role/1/", mock.Anything).Return(nil)
	client.On("DeleteJSON", mock.Anything, "configuration/v1/role/2/", mock.Anything).Return(nil)

	// Shared mock state - will be initialized on first create
	mockState := &config.LdapRole{}

	// Mock the CreateLdaprole API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ldap_role/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ldap_role/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.LdapRoleCreateRequest)
		// Reinitialize mockState from create request (important for destroy/recreate cycles)
		*mockState = config.LdapRole{
			ID:          123,
			ResourceURI: "/api/admin/configuration/v1/ldap_role/123/",
			Name:        createReq.Name,
			LdapGroupDN: createReq.LdapGroupDN,
			Roles:       createReq.Roles,
		}
	})

	// Mock the GetLdaprole API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ldap_role/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ldap_role := args.Get(3).(*config.LdapRole)
		*ldap_role = *mockState
	}).Maybe()

	// Mock the UpdateLdaprole API call
	client.On("PutJSON", mock.Anything, "configuration/v1/ldap_role/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.LdapRoleUpdateRequest)
		ldap_role := args.Get(3).(*config.LdapRole)

		// Update mock state based on request
		mockState.Name = updateRequest.Name
		mockState.LdapGroupDN = updateRequest.LdapGroupDN
		mockState.Roles = updateRequest.Roles

		// Return updated state
		*ldap_role = *mockState
	}).Maybe()

	// Mock the DeleteLdaprole API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ldap_role/123/"
	}), mock.Anything).Return(nil)

	testInfinityLdapRole(t, client)
}

func testInfinityLdapRole(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Test 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_role_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "name", "tf-test full"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "ldap_group_dn", "testfull"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "roles.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_ldap_role.ldap_role-test", "roles.*", "pexip_infinity_role.test1", "id"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_ldap_role.ldap_role-test", "roles.*", "pexip_infinity_role.test2", "id"),
				),
			},
			// Test 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_role_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "name", "tf-test min"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "ldap_group_dn", "testmin"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "roles.#", "0"),
				),
			},
			// Test 3: Destroy and recreate with min config
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_ldap_role_min"),
				Destroy: true,
			},
			// Test 4: Recreate with min config (after destroy)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_role_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "name", "tf-test min"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "ldap_group_dn", "testmin"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "roles.#", "0"),
				),
			},
			// Test 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_role_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "name", "tf-test full"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "ldap_group_dn", "testfull"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "roles.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_ldap_role.ldap_role-test", "roles.*", "pexip_infinity_role.test1", "id"),
					resource.TestCheckTypeSetElemAttrPair("pexip_infinity_ldap_role.ldap_role-test", "roles.*", "pexip_infinity_role.test2", "id"),
				),
			},
		},
	})
}
