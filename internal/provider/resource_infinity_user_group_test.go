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

func TestInfinityUserGroup(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking - starts with full config
	mockState := &config.UserGroup{
		ID:                      123,
		ResourceURI:             "/api/admin/configuration/v1/user_group/123/",
		Name:                    "tf-test-user-group-full",
		Description:             "tf-test user group description",
		Users:                   []string{},
		UserGroupEntityMappings: &[]config.UserGroupEntityMapping{},
	}

	// Mock the CreateUserGroup API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/user_group/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/user_group/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.UserGroupCreateRequest)
		// Update mock state based on create request
		mockState.Name = createReq.Name
		mockState.Description = createReq.Description
		if createReq.Users != nil {
			mockState.Users = createReq.Users
		} else {
			mockState.Users = []string{}
		}
		if createReq.UserGroupEntityMappings != nil {
			mappings := make([]config.UserGroupEntityMapping, len(createReq.UserGroupEntityMappings))
			for i, uri := range createReq.UserGroupEntityMappings {
				mappings[i] = config.UserGroupEntityMapping{ResourceURI: uri}
			}
			mockState.UserGroupEntityMappings = &mappings
		} else {
			mockState.UserGroupEntityMappings = &[]config.UserGroupEntityMapping{}
		}
	}).Maybe()

	// Mock the GetUserGroup API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/user_group/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		userGroup := args.Get(3).(*config.UserGroup)
		*userGroup = *mockState
	}).Maybe()

	// Mock the UpdateUserGroup API call
	client.On("PutJSON", mock.Anything, "configuration/v1/user_group/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.UserGroupUpdateRequest)
		userGroup := args.Get(3).(*config.UserGroup)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		// Description can be empty string, so always update
		mockState.Description = updateRequest.Description
		if updateRequest.Users != nil {
			mockState.Users = updateRequest.Users
		} else {
			mockState.Users = []string{}
		}
		if updateRequest.UserGroupEntityMappings != nil {
			// Convert string URIs to UserGroupEntityMapping objects
			mappings := make([]config.UserGroupEntityMapping, len(updateRequest.UserGroupEntityMappings))
			for i, uri := range updateRequest.UserGroupEntityMappings {
				mappings[i] = config.UserGroupEntityMapping{ResourceURI: uri}
			}
			mockState.UserGroupEntityMappings = &mappings
		} else {
			mockState.UserGroupEntityMappings = &[]config.UserGroupEntityMapping{}
		}

		// Return updated state
		*userGroup = *mockState
	}).Maybe()

	// Mock the DeleteUserGroup API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/user_group/123/"
	}), mock.Anything).Return(nil)

	testInfinityUserGroup(t, client)
}

func testInfinityUserGroup(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "name", "tf-test-user-group-full"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "description", "tf-test user group description"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "name", "tf-test-user-group"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "description", ""),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_user_group_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "name", "tf-test-user-group"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "description", ""),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.tf-test-user-group", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "name", "tf-test-user-group-full"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.tf-test-user-group", "description", "tf-test user group description"),
				),
			},
		},
	})
}
