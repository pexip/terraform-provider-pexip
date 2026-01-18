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

	// Mock the CreateLdaprole API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ldap_role/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ldap_role/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.LdapRole{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/ldap_role/123/",
		Name:        "ldap_role-test",
		LdapGroupDN: "test-value",
	}

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
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.LdapGroupDN != "" {
			mockState.LdapGroupDN = updateRequest.LdapGroupDN
		}

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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_role_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "name", "ldap_role-test"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_role_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_role.ldap_role-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_role.ldap_role-test", "name", "ldap_role-test"),
				),
			},
		},
	})
}
