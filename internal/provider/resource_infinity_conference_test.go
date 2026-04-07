/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"fmt"
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
	hostView := "two_mains_twentyone_pips"
	mockState := &config.Conference{
		ID:                              123,
		ResourceURI:                     "/api/admin/configuration/v1/conference/123/",
		Name:                            "tf-test-conference",
		AllowGuests:                     true,
		BreakoutRooms:                   true,
		CallType:                        "video-only",
		CryptoMode:                      test.StringPtr("besteffort"),
		DenoiseEnabled:                  true,
		Description:                     "Full test configuration for conference",
		DirectMedia:                     "best_effort",
		DirectMediaNotificationDuration: 5,
		EnableActiveSpeakerIndication:   true,
		EnableChat:                      "yes",
		EnableOverlayText:               true,
		ForcePresenterIntoMain:          true,
		GuestPIN:                        "654321",
		GuestsCanPresent:                true,
		GuestsCanSeeGuests:              "always",
		HostView:                        &hostView,
		LiveCaptionsEnabled:             "yes",
		MatchString:                     "^[0-9]+$",
		MaxCallRateIn:                   test.IntPtr(4096),
		MaxCallRateOut:                  test.IntPtr(2048),
		MaxPixelsPerSecond:              test.StringPtr("fullhd"),
		MuteAllGuests:                   true,
		NonIdpParticipants:              "allow_if_trusted",
		OnCompletion:                    test.StringPtr("{\"disconnect\": true}"),
		ParticipantLimit:                test.IntPtr(50),
		PostMatchString:                 "^test",
		PostReplaceString:               "new-test",
		PrimaryOwnerEmailAddress:        "owner@example.com",
		ReplaceString:                   "replaced",
		SoftmuteEnabled:                 true,
		SyncTag:                         "sync-123",
		TwoStageDialType:                "regular",
		ServiceType:                     "conference",
		PIN:                             "123456",
		Tag:                             "tf-test-tag",
		Aliases: &[]config.ConferenceAlias{
			{ID: 1},
			{ID: 2},
		},
		AutomaticParticipants: &[]config.AutomaticParticipant{
			{ID: 1},
			{ID: 2},
		},
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

		// Update all fields from the request
		mockState.Name = updateRequest.Name
		mockState.AllowGuests = updateRequest.AllowGuests
		mockState.BreakoutRooms = updateRequest.BreakoutRooms
		mockState.CallType = updateRequest.CallType
		mockState.CryptoMode = updateRequest.CryptoMode
		mockState.DenoiseEnabled = updateRequest.DenoiseEnabled
		mockState.Description = updateRequest.Description
		mockState.DirectMedia = updateRequest.DirectMedia
		mockState.DirectMediaNotificationDuration = updateRequest.DirectMediaNotificationDuration
		mockState.EnableActiveSpeakerIndication = updateRequest.EnableActiveSpeakerIndication
		mockState.EnableChat = updateRequest.EnableChat
		mockState.EnableOverlayText = updateRequest.EnableOverlayText
		mockState.ForcePresenterIntoMain = updateRequest.ForcePresenterIntoMain
		mockState.GuestPIN = updateRequest.GuestPIN
		mockState.GuestsCanPresent = updateRequest.GuestsCanPresent
		mockState.GuestsCanSeeGuests = updateRequest.GuestsCanSeeGuests
		mockState.LiveCaptionsEnabled = updateRequest.LiveCaptionsEnabled
		mockState.MatchString = updateRequest.MatchString
		mockState.MuteAllGuests = updateRequest.MuteAllGuests
		mockState.NonIdpParticipants = updateRequest.NonIdpParticipants
		mockState.PIN = updateRequest.PIN
		mockState.PostMatchString = updateRequest.PostMatchString
		mockState.PostReplaceString = updateRequest.PostReplaceString
		mockState.PrimaryOwnerEmailAddress = updateRequest.PrimaryOwnerEmailAddress
		mockState.ReplaceString = updateRequest.ReplaceString
		mockState.ServiceType = updateRequest.ServiceType
		mockState.SoftmuteEnabled = updateRequest.SoftmuteEnabled
		mockState.SyncTag = updateRequest.SyncTag
		mockState.Tag = updateRequest.Tag
		mockState.TwoStageDialType = updateRequest.TwoStageDialType

		// Handle pointer fields
		mockState.GMSAccessToken = updateRequest.GMSAccessToken
		mockState.GuestIdentityProviderGroup = updateRequest.GuestIdentityProviderGroup
		mockState.GuestView = updateRequest.GuestView
		mockState.HostIdentityProviderGroup = updateRequest.HostIdentityProviderGroup
		mockState.HostView = updateRequest.HostView
		mockState.MaxCallRateIn = updateRequest.MaxCallRateIn
		mockState.MaxCallRateOut = updateRequest.MaxCallRateOut
		mockState.MaxPixelsPerSecond = updateRequest.MaxPixelsPerSecond
		mockState.MediaPlaylist = updateRequest.MediaPlaylist
		mockState.MSSIPProxy = updateRequest.MSSIPProxy
		mockState.OnCompletion = updateRequest.OnCompletion
		mockState.ParticipantLimit = updateRequest.ParticipantLimit
		mockState.PinningConfig = updateRequest.PinningConfig
		mockState.SystemLocation = updateRequest.SystemLocation
		mockState.TeamsProxy = updateRequest.TeamsProxy

		// IVRTheme is a special case - updateRequest has *string (URI), but mockState has *IVRTheme object
		// For mock purposes, we can leave it as nil since we're not testing IVRTheme
		if updateRequest.IVRTheme != nil {
			mockState.IVRTheme = nil // In real API, this would be resolved to an object
		} else {
			mockState.IVRTheme = nil
		}

		// Aliases - if provided, update; if nil/empty, clear
		if updateRequest.Aliases != nil {
			var aliasObjects []config.ConferenceAlias
			for i, alias := range *updateRequest.Aliases {
				aliasObjects = append(aliasObjects, config.ConferenceAlias{
					ID:          i + 1,
					Alias:       alias,
					Conference:  "/api/admin/configuration/v1/conference/123/",
					ResourceURI: fmt.Sprintf("/api/admin/configuration/v1/conference_alias/%d/", i+1),
				})
			}
			mockState.Aliases = &aliasObjects
		} else {
			mockState.Aliases = nil
		}
		// AutomaticParticipants - if provided, update; if nil/empty, clear
		if updateRequest.AutomaticParticipants != nil {
			var participantObjects []config.AutomaticParticipant
			for i, participant := range updateRequest.AutomaticParticipants {
				// Do not set ResourceURI: the real API does not populate it on embedded objects.
				participantObjects = append(participantObjects, config.AutomaticParticipant{
					ID:                  i + 1,
					Alias:               participant,
					Protocol:            "auto",
					CallType:            "video",
					Role:                "guest",
					KeepConferenceAlive: "keep_conference_alive_never",
					Routing:             "routing_rule",
				})
			}
			mockState.AutomaticParticipants = &participantObjects
		} else {
			mockState.AutomaticParticipants = nil
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
			// Test 1: Create with full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "description", "Full test configuration for conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "service_type", "conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "pin", "123456"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guest_pin", "654321"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "tag", "tf-test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "allow_guests", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "breakout_rooms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "call_type", "video-only"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "crypto_mode", "besteffort"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "denoise_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "direct_media", "best_effort"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "direct_media_notification_duration", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "enable_active_speaker_indication", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "enable_chat", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "enable_overlay_text", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "force_presenter_into_main", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guests_can_present", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guests_can_see_guests", "always"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "host_view", "two_mains_twentyone_pips"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "live_captions_enabled", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "match_string", "^[0-9]+$"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "max_callrate_in", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "max_callrate_out", "2048"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "mute_all_guests", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "non_idp_participants", "allow_if_trusted"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "on_completion", "{\"disconnect\": true}"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "participant_limit", "50"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "post_match_string", "^test"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "post_replace_string", "new-test"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "primary_owner_email_address", "owner@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "replace_string", "replaced"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "softmute_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "sync_tag", "sync-123"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "two_stage_dial_type", "regular"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "automatic_participants.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "automatic_participants.*", "/api/admin/configuration/v1/automatic_participant/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "automatic_participants.*", "/api/admin/configuration/v1/automatic_participant/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "aliases.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "aliases.*", "/api/admin/configuration/v1/conference_alias/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "aliases.*", "/api/admin/configuration/v1/conference_alias/2/"),
				),
			},
			// Test 2: Update to min configuration, then delete
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "service_type", "conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "pin", ""),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guest_pin", ""),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "allow_guests", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "tag", ""),
				),
			},
			// Test 3: Create with min configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "service_type", "conference"),
				),
			},
			// Test 4: Update to full configuration
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference.tf-test-conference", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "name", "tf-test-conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "description", "Full test configuration for conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "service_type", "conference"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "pin", "123456"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guest_pin", "654321"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "tag", "tf-test-tag"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "allow_guests", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "breakout_rooms", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "call_type", "video-only"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "crypto_mode", "besteffort"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "denoise_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "direct_media", "best_effort"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "direct_media_notification_duration", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "enable_active_speaker_indication", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "enable_chat", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "enable_overlay_text", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "force_presenter_into_main", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guests_can_present", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "guests_can_see_guests", "always"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "host_view", "two_mains_twentyone_pips"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "live_captions_enabled", "yes"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "match_string", "^[0-9]+$"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "max_callrate_in", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "max_callrate_out", "2048"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "max_pixels_per_second", "fullhd"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "mute_all_guests", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "non_idp_participants", "allow_if_trusted"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "on_completion", "{\"disconnect\": true}"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "participant_limit", "50"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "post_match_string", "^test"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "post_replace_string", "new-test"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "primary_owner_email_address", "owner@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "replace_string", "replaced"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "softmute_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "sync_tag", "sync-123"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "two_stage_dial_type", "regular"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "automatic_participants.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "automatic_participants.*", "/api/admin/configuration/v1/automatic_participant/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "automatic_participants.*", "/api/admin/configuration/v1/automatic_participant/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_conference.tf-test-conference", "aliases.#", "2"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "aliases.*", "/api/admin/configuration/v1/conference_alias/1/"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_conference.tf-test-conference", "aliases.*", "/api/admin/configuration/v1/conference_alias/2/"),
				),
			},
		},
	})
}
