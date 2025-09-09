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

func TestInfinityLdapSyncField(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateLdapsyncfield API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ldap_sync_field/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ldap_sync_field/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.LdapSyncField{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/ldap_sync_field/123/",
		Name:                 "ldap_sync_field-test",
		Description:          "Test LdapSyncField",
		TemplateVariableName: "ldap_sync_field-test",
		IsBinary:             true,
	}

	// Mock the GetLdapsyncfield API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/ldap_sync_field/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ldap_sync_field := args.Get(3).(*config.LdapSyncField)
		*ldap_sync_field = *mockState
	}).Maybe()

	// Mock the UpdateLdapsyncfield API call
	client.On("PutJSON", mock.Anything, "configuration/v1/ldap_sync_field/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.LdapSyncFieldUpdateRequest)
		ldap_sync_field := args.Get(3).(*config.LdapSyncField)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.TemplateVariableName != "" {
			mockState.TemplateVariableName = updateRequest.TemplateVariableName
		}
		if updateRequest.IsBinary != nil {
			mockState.IsBinary = *updateRequest.IsBinary
		}

		// Return updated state
		*ldap_sync_field = *mockState
	}).Maybe()

	// Mock the DeleteLdapsyncfield API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/ldap_sync_field/123/"
	}), mock.Anything).Return(nil)

	testInfinityLdapSyncField(t, client)
}

func testInfinityLdapSyncField(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "name", "ldap_sync_field-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "description", "Test LdapSyncField"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "template_variable_name", "ldap_sync_field-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "is_binary", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "name", "ldap_sync_field-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "description", "Updated Test LdapSyncField"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "template_variable_name", "ldap_sync_field-test"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.ldap_sync_field-test", "is_binary", "false"),
				),
			},
		},
	})
}
