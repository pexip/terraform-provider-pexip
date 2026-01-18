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

func TestInfinityConference(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateConference API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/conference/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/conference/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	hostView := "one_main_seven_pips"
	mockState := &config.Conference{
		ID:            123,
		ResourceURI:   "/api/admin/configuration/v1/conference/123/",
		Name:          "conference-test",
		AllowGuests:   false,
		BreakoutRooms: false,
		CallType:      "video",
		//CreationTime:  "2024-01-01T12:00:00Z", --- IGNORE ---
		DenoiseEnabled:                  false,
		Description:                     "Test Conference",
		DirectMedia:                     "never",
		DirectMediaNotificationDuration: 0,
		EnableActiveSpeakerIndication:   false,
		EnableChat:                      "default",
		EnableOverlayText:               false,
		ForcePresenterIntoMain:          false,
		GuestPIN:                        "",
		GuestsCanPresent:                true,
		GuestsCanSeeGuests:              "no_hosts",
		GuestView:                       nil,
		HostView:                        &hostView,
		LiveCaptionsEnabled:             "default",
		MatchString:                     "",
		MuteAllGuests:                   false,
		NonIdpParticipants:              "disallow_all",
		PostMatchString:                 "",
		PrimaryOwnerEmailAddress:        "",
		ReplaceString:                   "",
		SoftmuteEnabled:                 false,
		SyncTag:                         "",
		TwoStageDialType:                "regular",

		ServiceType: "conference",
		PIN:         "1234",
		Tag:         "test-tag",
	}

	// Mock the GetConference API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/conference/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		conference := args.Get(3).(*config.Conference)
		*conference = *mockState
	}).Maybe()

	// Mock the UpdateConference API call
	client.On("PutJSON", mock.Anything, "configuration/v1/conference/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ConferenceUpdateRequest)
		conference := args.Get(3).(*config.Conference)

		// Update mock state
		mockState.Name = updateRequest.Name
		mockState.AllowGuests = true
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.PIN != "" {
			mockState.PIN = updateRequest.PIN
		}
		if updateRequest.GuestPIN != "" {
			mockState.GuestPIN = updateRequest.GuestPIN
		}
		if updateRequest.Tag != "" {
			mockState.Tag = updateRequest.Tag
		}
		// Return updated state
		*conference = *mockState
	}).Maybe()

	// Mock the DeleteConference API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/conference/123/"
	}), mock.Anything).Return(nil)

	testInfinityConference(t, client)
}

func testInfinityConference(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.conference-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.conference-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "name", "conference-test"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "description", "Test Conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "service_type", "conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "pin", "1234"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "tag", "test-tag"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.conference-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.conference-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "name", "conference-test"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "description", "Updated Test Conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "service_type", "conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "pin", "9876"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "guest_pin", "4321"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "allow_guests", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.conference-test", "tag", "updated-tag"),
				),
			},
		},
	})
}
