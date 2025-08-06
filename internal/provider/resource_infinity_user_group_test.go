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

func TestInfinityUserGroup(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateUserGroup API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/user_group/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/user_group/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.UserGroup{
		ID:                      123,
		ResourceURI:             "/api/admin/configuration/v1/user_group/123/",
		Name:                    "user-group-test",
		Description:             "Test User Group",
		Users:                   []string{},
		UserGroupEntityMappings: []string{},
	}

	// Mock the GetUserGroup API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/user_group/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		userGroup := args.Get(2).(*config.UserGroup)
		*userGroup = *mockState
	}).Maybe()

	// Mock the UpdateUserGroup API call
	client.On("PutJSON", mock.Anything, "configuration/v1/user_group/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.UserGroupUpdateRequest)
		userGroup := args.Get(3).(*config.UserGroup)

		// Update mock state
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.Users != nil {
			mockState.Users = updateRequest.Users
		}
		if updateRequest.UserGroupEntityMappings != nil {
			mockState.UserGroupEntityMappings = updateRequest.UserGroupEntityMappings
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
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.user-group-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.user-group-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.user-group-test", "name", "user-group-test"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.user-group-test", "description", "Test User Group"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.user-group-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group.user-group-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.user-group-test", "name", "user-group-test"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group.user-group-test", "description", "Updated Test User Group"),
				),
			},
		},
	})
}
