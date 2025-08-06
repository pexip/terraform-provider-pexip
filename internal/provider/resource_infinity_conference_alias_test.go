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

func TestInfinityConferenceAlias(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

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
		Alias:       "test-value",
		Description: "Test ConferenceAlias",
		Conference:  "test-value",
	}

	// Mock the GetConferencealias API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/conference_alias/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		conference_alias := args.Get(2).(*config.ConferenceAlias)
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
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_alias_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.conference_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.conference_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.conference_alias-test", "description", "Test ConferenceAlias"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_conference_alias_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.conference_alias-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_conference_alias.conference_alias-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_conference_alias.conference_alias-test", "description", "Updated Test ConferenceAlias"),
				),
			},
		},
	})
}
