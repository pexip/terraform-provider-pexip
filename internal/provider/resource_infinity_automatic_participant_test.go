/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
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

	client := infinity.NewClientMock()

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
		Conference:          "test-conference",
		Protocol:            "sip",
		CallType:            "video",
		Role:                "guest",
		DTMFSequence:        "123#",
		KeepConferenceAlive: "keep_conference_alive",
		Routing:             "auto",
		SystemLocation:      test.StringPtr("test-location"),
		Streaming:           true,
		RemoteDisplayName:   "automatic_participant-test",
		PresentationURL:     "https://example.com",
	}

	client.On("GetJSON", mock.Anything, "configuration/v1/automatic_participant/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		participant := args.Get(2).(*config.AutomaticParticipant)
		*participant = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/automatic_participant/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.AutomaticParticipantUpdateRequest)
		participant := args.Get(3).(*config.AutomaticParticipant)
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
				),
			},
		},
	})
}
