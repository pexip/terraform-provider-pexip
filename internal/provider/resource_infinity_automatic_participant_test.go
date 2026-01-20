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

func TestInfinityAutomaticParticipant(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

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
		Name:                            "test-conference",
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

	// Mock system_location creation
	locationCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/system_location/1/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/system_location/", mock.Anything, mock.Anything).Return(locationCreateResponse, nil)

	locationState := &config.SystemLocation{
		ID:                        1,
		ResourceURI:               "/api/admin/configuration/v1/system_location/1/",
		Name:                      "test-location",
		Description:               "Test Location",
		MTU:                       1460,
		SignallingQoS:             test.IntPtr(0),
		MediaQoS:                  test.IntPtr(0),
		BDPMPinChecksEnabled:      "GLOBAL",
		BDPMScanQuarantineEnabled: "GLOBAL",
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/system_location/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		loc := args.Get(3).(*config.SystemLocation)
		*loc = *locationState
	}).Maybe()

	client.On("PatchJSON", mock.Anything, "configuration/v1/system_location/1/", mock.Anything, mock.Anything).Return(nil).Maybe()
	client.On("DeleteJSON", mock.Anything, "configuration/v1/system_location/1/", mock.Anything).Return(nil).Maybe()

	// Mock automatic_participant creation
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/automatic_participant/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/automatic_participant/", mock.Anything, mock.Anything).Return(createResponse, nil)

	mockState := &config.AutomaticParticipant{
		ID:                  123,
		ResourceURI:         "/api/admin/configuration/v1/automatic_participant/123/",
		Alias:               "automatic-participant-test",
		Description:         "Test AutomaticParticipant",
		Conference:          []string{"/api/admin/configuration/v1/conference/1/"},
		Protocol:            "sip",
		CallType:            "video",
		Role:                "guest",
		DTMFSequence:        "123#",
		KeepConferenceAlive: "keep_conference_alive",
		Routing:             "manual",
		SystemLocation:      test.StringPtr("/api/admin/configuration/v1/system_location/1/"),
		Streaming:           true,
		RemoteDisplayName:   "automatic_participant-test",
		PresentationURL:     "https://example.com",
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/automatic_participant/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		participant := args.Get(3).(*config.AutomaticParticipant)
		*participant = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/automatic_participant/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.AutomaticParticipantUpdateRequest)
		participant := args.Get(3).(*config.AutomaticParticipant)

		// Update the mock state with all fields from the update request
		mockState.Alias = updateRequest.Alias
		mockState.Description = updateRequest.Description
		mockState.Conference = updateRequest.Conference
		mockState.Protocol = updateRequest.Protocol
		mockState.CallType = updateRequest.CallType
		mockState.Role = updateRequest.Role
		mockState.DTMFSequence = updateRequest.DTMFSequence
		mockState.KeepConferenceAlive = updateRequest.KeepConferenceAlive
		mockState.Routing = updateRequest.Routing
		mockState.SystemLocation = updateRequest.SystemLocation
		if updateRequest.Streaming != nil {
			mockState.Streaming = *updateRequest.Streaming
		}
		mockState.RemoteDisplayName = updateRequest.RemoteDisplayName
		mockState.PresentationURL = updateRequest.PresentationURL

		// Return the updated state
		*participant = *mockState
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/automatic_participant/123/"
	}), mock.Anything).Return(nil)

	testInfinityAutomaticParticipant(t, client)
}

func testInfinityAutomaticParticipant(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_automatic_participant_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_automatic_participant.automatic-participant-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_automatic_participant.automatic-participant-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_automatic_participant.automatic-participant-test", "alias", "automatic-participant-test"),
					resource.TestCheckResourceAttr("pexip_infinity_automatic_participant.automatic-participant-test", "description", "Test AutomaticParticipant"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_automatic_participant_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_automatic_participant.automatic-participant-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_automatic_participant.automatic-participant-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_automatic_participant.automatic-participant-test", "alias", "automatic-participant-updated"),
					resource.TestCheckResourceAttr("pexip_infinity_automatic_participant.automatic-participant-test", "description", "Updated Test AutomaticParticipant"),
				),
			},
		},
	})
}
