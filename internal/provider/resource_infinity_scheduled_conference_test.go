/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"
	"time"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/pexip/go-infinity-sdk/v38/util"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityScheduledConference(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateScheduledconference API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/scheduled_conference/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/scheduled_conference/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	startTime, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-01-01T11:00:00Z")

	mockState := &config.ScheduledConference{
		ID:                  123,
		ResourceURI:         "/api/admin/configuration/v1/scheduled_conference/123/",
		Conference:          "test-value",
		StartTime:           util.InfinityTime{Time: startTime},
		EndTime:             util.InfinityTime{Time: endTime},
		Subject:             "test-value",
		EWSItemID:           "test-value",
		EWSItemUID:          "test-value",
		RecurringConference: test.StringPtr("test-value"),
		ScheduledAlias:      test.StringPtr("test-value"),
	}

	// Mock the GetScheduledconference API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/scheduled_conference/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		scheduled_conference := args.Get(3).(*config.ScheduledConference)
		*scheduled_conference = *mockState
	}).Maybe()

	// Mock the UpdateScheduledconference API call
	client.On("PutJSON", mock.Anything, "configuration/v1/scheduled_conference/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ScheduledConferenceUpdateRequest)
		scheduled_conference := args.Get(3).(*config.ScheduledConference)

		// Update mock state based on request
		if updateRequest.Conference != "" {
			mockState.Conference = updateRequest.Conference
		}
		if updateRequest.Subject != "" {
			mockState.Subject = updateRequest.Subject
		}
		if updateRequest.EWSItemID != "" {
			mockState.EWSItemID = updateRequest.EWSItemID
		}
		if updateRequest.EWSItemUID != "" {
			mockState.EWSItemUID = updateRequest.EWSItemUID
		}
		if updateRequest.RecurringConference != nil {
			mockState.RecurringConference = updateRequest.RecurringConference
		}
		if updateRequest.ScheduledAlias != nil {
			mockState.ScheduledAlias = updateRequest.ScheduledAlias
		}
		if updateRequest.StartTime != nil && !updateRequest.StartTime.IsZero() {
			mockState.StartTime = *updateRequest.StartTime
		}
		if updateRequest.EndTime != nil && !updateRequest.EndTime.IsZero() {
			mockState.EndTime = *updateRequest.EndTime
		}

		// Return updated state
		*scheduled_conference = *mockState
	}).Maybe()

	// Mock the DeleteScheduledconference API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/scheduled_conference/123/"
	}), mock.Anything).Return(nil)

	testInfinityScheduledConference(t, client)
}

func testInfinityScheduledConference(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_conference_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_conference.scheduled_conference-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_conference.scheduled_conference-test", "resource_id"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_conference_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_conference.scheduled_conference-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_conference.scheduled_conference-test", "resource_id"),
				),
			},
		},
	})
}
