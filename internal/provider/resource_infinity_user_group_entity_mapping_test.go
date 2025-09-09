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

func TestInfinityUserGroupEntityMapping(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateUsergroupentitymapping API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/user_group_entity_mapping/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/user_group_entity_mapping/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.UserGroupEntityMapping{
		ID:                123,
		ResourceURI:       "/api/admin/configuration/v1/user_group_entity_mapping/123/",
		Description:       "Test UserGroupEntityMapping",
		EntityResourceURI: "test-value",
		UserGroup:         "test-value",
	}

	// Mock the GetUsergroupentitymapping API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/user_group_entity_mapping/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		user_group_entity_mapping := args.Get(3).(*config.UserGroupEntityMapping)
		*user_group_entity_mapping = *mockState
	}).Maybe()

	// Mock the UpdateUsergroupentitymapping API call
	client.On("PutJSON", mock.Anything, "configuration/v1/user_group_entity_mapping/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.UserGroupEntityMappingUpdateRequest)
		user_group_entity_mapping := args.Get(3).(*config.UserGroupEntityMapping)

		// Update mock state based on request
		if updateReq.Description != "" {
			mockState.Description = updateReq.Description
		}
		if updateReq.EntityResourceURI != "" {
			mockState.EntityResourceURI = updateReq.EntityResourceURI
		}
		if updateReq.UserGroup != "" {
			mockState.UserGroup = updateReq.UserGroup
		}

		// Return updated state
		*user_group_entity_mapping = *mockState
	}).Maybe()

	// Mock the DeleteUsergroupentitymapping API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/user_group_entity_mapping/123/"
	}), mock.Anything).Return(nil)

	testInfinityUserGroupEntityMapping(t, client)
}

func testInfinityUserGroupEntityMapping(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "description", "Test UserGroupEntityMapping"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "entity_resource_uri", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "user_group", "test-value"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "description", "Updated Test UserGroupEntityMapping"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "entity_resource_uri", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.user_group_entity_mapping-test", "user_group", "updated-value"),
				),
			},
		},
	})
}
