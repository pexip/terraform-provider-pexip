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

func TestInfinityLdapSyncField(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared mock state - will be initialized on first create
	mockState := &config.LdapSyncField{}

	// Mock the CreateLdapsyncfield API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/ldap_sync_field/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/ldap_sync_field/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.LdapSyncFieldCreateRequest)
		// Reinitialize mockState from create request (important for destroy/recreate cycles)
		*mockState = config.LdapSyncField{
			ID:                   123,
			ResourceURI:          "/api/admin/configuration/v1/ldap_sync_field/123/",
			Name:                 createReq.Name,
			Description:          createReq.Description,
			TemplateVariableName: createReq.TemplateVariableName,
			IsBinary:             createReq.IsBinary,
		}
	})

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
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description
		mockState.TemplateVariableName = updateRequest.TemplateVariableName
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
			// Test 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "name", "tf-test-full"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "template_variable_name", "testfull"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "is_binary", "true"),
				),
			},
			// Test 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "name", "tf-test-min"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "template_variable_name", "testmin"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "is_binary", "false"),
				),
			},
			// Test 3: Destroy and recreate with min config
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_min"),
				Destroy: true,
			},
			// Test 4: Recreate with min config (after destroy)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "name", "tf-test-min"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "template_variable_name", "testmin"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "is_binary", "false"),
				),
			},
			// Test 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_ldap_sync_field_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_ldap_sync_field.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "name", "tf-test-full"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "template_variable_name", "testfull"),
					resource.TestCheckResourceAttr("pexip_infinity_ldap_sync_field.test", "is_binary", "true"),
				),
			},
		},
	})
}
