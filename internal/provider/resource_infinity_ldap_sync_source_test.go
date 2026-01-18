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

func TestInfinityLdapSyncSource(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateLdapsyncsource API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ldap_sync_source/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ldap_sync_source/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.LdapSyncSource{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/ldap_sync_source/123/",
		Name:                 "ldap_sync_source-test",
		Description:          "Test LdapSyncSource",
		LdapServer:           "test-value",
		LdapBaseDN:           "test-value",
		LdapBindUsername:     "ldap_sync_source-test",
		LdapBindPassword:     "test-value",
		LdapUseGlobalCatalog: true,
		LdapPermitNoTLS:      true,
	}

	// Mock the GetLdapsyncsource API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ldap_sync_source/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ldap_sync_source := args.Get(3).(*config.LdapSyncSource)
		*ldap_sync_source = *mockState
	}).Maybe()

	// Mock the UpdateLdapsyncsource API call
	client.On("PutJSON", mock.Anything, "configuration/v1/ldap_sync_source/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.LdapSyncSourceUpdateRequest)
		ldap_sync_source := args.Get(3).(*config.LdapSyncSource)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.LdapServer != "" {
			mockState.LdapServer = updateRequest.LdapServer
		}
		if updateRequest.LdapBaseDN != "" {
			mockState.LdapBaseDN = updateRequest.LdapBaseDN
		}
		if updateRequest.LdapBindUsername != "" {
			mockState.LdapBindUsername = updateRequest.LdapBindUsername
		}
		if updateRequest.LdapBindPassword != "" {
			mockState.LdapBindPassword = updateRequest.LdapBindPassword
		}
		if updateRequest.LdapUseGlobalCatalog != nil {
			mockState.LdapUseGlobalCatalog = *updateRequest.LdapUseGlobalCatalog
		}
		if updateRequest.LdapPermitNoTLS != nil {
			mockState.LdapPermitNoTLS = *updateRequest.LdapPermitNoTLS
		}

		// Return updated state
		*ldap_sync_source = *mockState
	}).Maybe()

	// Mock the DeleteLdapsyncsource API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ldap_sync_source/123/"
	}), mock.Anything).Return(nil)

	testInfinityLdapSyncSource(t, client)
}

func testInfinityLdapSyncSource(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_source_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "name", "ldap_sync_source-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "description", "Test LdapSyncSource"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "ldap_bind_username", "ldap_sync_source-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "ldap_use_global_catalog", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "ldap_permit_no_tls", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_source_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "name", "ldap_sync_source-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "description", "Updated Test LdapSyncSource"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "ldap_bind_username", "ldap_sync_source-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "ldap_use_global_catalog", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_source.ldap_sync_source-test", "ldap_permit_no_tls", "false"),
				),
			},
		},
	})
}
