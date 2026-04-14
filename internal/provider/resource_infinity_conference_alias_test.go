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

func TestInfinityConferenceAlias(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

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
		Name:                            "tf-test-conference",
		Description:                     "Test Conference",
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
		HostView:                        &hostView,
		LiveCaptionsEnabled:             "default",
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
		conf := args.Get(3).(*config.Conference)
		*conf = *conferenceState
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/conference/1/", mock.Anything).Return(nil).Maybe()

	// Mock the CreateConferencealias API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/conference_alias/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/conference_alias/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ConferenceAlias{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/conference_alias/123/",
		Alias:       "tf-test-alias",
		Description: "Test Conference Alias Description",
		Conference:  "/api/admin/configuration/v1/conference/1/",
	}

	// Mock the GetConferencealias API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/conference_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		conference_alias := args.Get(3).(*config.ConferenceAlias)
		*conference_alias = *mockState
	}).Maybe()

	// Mock the UpdateConferencealias API call
	client.On("PutJSON", mock.Anything, "configuration/v1/conference_alias/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ConferenceAliasUpdateRequest)
		conference_alias := args.Get(3).(*config.ConferenceAlias)

		// Update mock state based on request
		mockState.Alias = updateRequest.Alias
		mockState.Description = updateRequest.Description
		mockState.Conference = updateRequest.Conference

		// Return updated state
		*conference_alias = *mockState
	}).Maybe()

	// Mock the DeleteConferencealias API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/conference_alias/123/"
	}), mock.Anything).Return(nil)

	testInfinityConferenceAlias(t, client)
}

func testInfinityConferenceAlias(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_alias_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.tf-test-conference-alias", "alias", "tf-test-alias"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.tf-test-conference-alias", "description", "Test Conference Alias Description"),
				),
			},
			// Step 2: Update to min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_alias_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.tf-test-conference-alias", "alias", "tf-test-alias"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_conference_alias_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_alias_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.tf-test-conference-alias", "alias", "tf-test-alias"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_alias_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.tf-test-conference-alias", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.tf-test-conference-alias", "alias", "tf-test-alias"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.tf-test-conference-alias", "description", "Test Conference Alias Description"),
				),
			},
		},
	})
}
