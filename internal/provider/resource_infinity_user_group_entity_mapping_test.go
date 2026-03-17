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

func TestInfinityUserGroupEntityMapping(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock user group creation
	userGroupCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/user_group/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/user_group/", mock.Anything, mock.Anything).Return(userGroupCreateResponse, nil)

	userGroupState := &config.UserGroup{
		ID:                      1,
		ResourceURI:             "/api/admin/configuration/v1/user_group/1/",
		Name:                    "tf-test-user-group-full",
		Description:             "",
		Users:                   []string{},
		UserGroupEntityMappings: &[]config.UserGroupEntityMapping{},
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/user_group/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		userGroup := args.Get(3).(*config.UserGroup)
		*userGroup = *userGroupState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/user_group/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.UserGroupUpdateRequest)
		userGroup := args.Get(3).(*config.UserGroup)
		if updateReq.Name != "" {
			userGroupState.Name = updateReq.Name
		}
		userGroupState.Description = updateReq.Description
		*userGroup = *userGroupState
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/user_group/1/", mock.Anything).Return(nil)

	// Mock conference creation
	conferenceCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/conference/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/conference/", mock.Anything, mock.Anything).Return(conferenceCreateResponse, nil)

	hostView := "one_main_seven_pips"
	conferenceState := &config.Conference{
		ID:                              1,
		ResourceURI:                     "/api/admin/configuration/v1/conference/1/",
		Name:                            "tf-test-conference-full",
		Description:                     "",
		ServiceType:                     "conference",
		AllowGuests:                     false,
		BreakoutRooms:                   false,
		CallType:                        "video",
		CryptoMode:                      test.StringPtr(""), // API returns empty string, not nil
		DenoiseEnabled:                  false,
		DirectMedia:                     "never",
		DirectMediaNotificationDuration: 0,
		EnableActiveSpeakerIndication:   false,
		EnableChat:                      "default",
		EnableOverlayText:               false,
		ForcePresenterIntoMain:          false,
		GuestPIN:                        "",
		GuestsCanPresent:                true,
		GuestsCanSeeGuests:              "no_hosts",
		LiveCaptionsEnabled:             "default",
		HostView:                        &hostView,
		MatchString:                     "",
		MuteAllGuests:                   false,
		NonIdpParticipants:              "disallow_all",
		OnCompletion:                    test.StringPtr(""), // API returns empty string, not nil
		PostMatchString:                 "",
		PostReplaceString:               "",
		PrimaryOwnerEmailAddress:        "",
		ReplaceString:                   "",
		SoftmuteEnabled:                 false,
		SyncTag:                         "",
		Tag:                             "",
		TwoStageDialType:                "regular",
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/conference/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		conference := args.Get(3).(*config.Conference)
		*conference = *conferenceState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/conference/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.ConferenceUpdateRequest)
		conference := args.Get(3).(*config.Conference)
		if updateReq.Name != "" {
			conferenceState.Name = updateReq.Name
		}
		conferenceState.Description = updateReq.Description
		*conference = *conferenceState
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/conference/1/", mock.Anything).Return(nil)

	// Shared state for mocking - starts with full config
	mockState := &config.UserGroupEntityMapping{
		ID:                123,
		ResourceURI:       "/api/admin/configuration/v1/user_group_entity_mapping/123/",
		Description:       "tf-test user group entity mapping description",
		EntityResourceURI: "/api/admin/configuration/v1/conference/1/",
		UserGroup:         "/api/admin/configuration/v1/user_group/1/",
	}

	// Mock the CreateUsergroupentitymapping API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/user_group_entity_mapping/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/user_group_entity_mapping/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.UserGroupEntityMappingCreateRequest)
		mockState.Description = createReq.Description
		mockState.EntityResourceURI = createReq.EntityResourceURI
		mockState.UserGroup = createReq.UserGroup
	}).Maybe()

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
		mockState.Description = updateReq.Description
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
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "description", "tf-test user group entity mapping description"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "entity_resource_uri"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "user_group"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "description", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "entity_resource_uri"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "user_group"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "description", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "entity_resource_uri"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "user_group"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_user_group_entity_mapping_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "description", "tf-test user group entity mapping description"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "entity_resource_uri"),
					resource.TestCheckResourceAttrSet("pexip_infinity_user_group_entity_mapping.tf-test-mapping", "user_group"),
				),
			},
		},
	})
}
