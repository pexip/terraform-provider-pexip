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

func TestInfinityRecurringConference(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateRecurringconference API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/recurring_conference/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/recurring_conference/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.RecurringConference{
		ID:             123,
		ResourceURI:    "/api/admin/configuration/v1/recurring_conference/123/",
		Conference:     "test-value",
		CurrentIndex:   1,
		EWSItemID:      "test-value",
		IsDepleted:     true,
		Subject:        "test-value",
		ScheduledAlias: test.StringPtr("test-value"),
	}

	// Mock the GetRecurringconference API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/recurring_conference/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		recurring_conference := args.Get(2).(*config.RecurringConference)
		*recurring_conference = *mockState
	}).Maybe()

	// Mock the UpdateRecurringconference API call
	client.On("PutJSON", mock.Anything, "configuration/v1/recurring_conference/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.RecurringConferenceUpdateRequest)
		recurring_conference := args.Get(3).(*config.RecurringConference)

		// Update mock state based on request
		if updateRequest.Conference != "" {
			mockState.Conference = updateRequest.Conference
		}
		if updateRequest.CurrentIndex != nil {
			mockState.CurrentIndex = *updateRequest.CurrentIndex
		}
		if updateRequest.EWSItemID != "" {
			mockState.EWSItemID = updateRequest.EWSItemID
		}
		if updateRequest.IsDepleted != nil {
			mockState.IsDepleted = *updateRequest.IsDepleted
		}
		if updateRequest.Subject != "" {
			mockState.Subject = updateRequest.Subject
		}
		if updateRequest.ScheduledAlias != nil {
			mockState.ScheduledAlias = updateRequest.ScheduledAlias
		}

		// Return updated state
		*recurring_conference = *mockState
	}).Maybe()

	// Mock the DeleteRecurringconference API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/recurring_conference/123/"
	}), mock.Anything).Return(nil)

	testInfinityRecurringConference(t, client)
}

func testInfinityRecurringConference(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_recurring_conference_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_recurring_conference.recurring_conference-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_recurring_conference.recurring_conference-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_recurring_conference.recurring_conference-test", "is_depleted", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_recurring_conference_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_recurring_conference.recurring_conference-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_recurring_conference.recurring_conference-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_recurring_conference.recurring_conference-test", "is_depleted", "false"),
				),
			},
		},
	})
}
