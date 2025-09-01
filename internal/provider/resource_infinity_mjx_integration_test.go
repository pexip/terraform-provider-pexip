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

func TestInfinityMjxIntegration(t *testing.T) {

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateMjxintegration API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_integration/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_integration/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.MjxIntegration{
		ID:                          123,
		ResourceURI:                 "/api/admin/configuration/v1/mjx_integration/123/",
		Name:                        "mjx_integration-test",
		Description:                 "Test MjxIntegration",
		StartBuffer:                 30,
		EndBuffer:                   30,
		DisplayUpcomingMeetings:     12,
		EnableNonVideoMeetings:      true,
		EnablePrivateMeetings:       true,
		EPUsername:                  "mjx_integration-test",
		EPPassword:                  "test-value",
		EPUseHTTPS:                  true,
		EPVerifyCertificate:         true,
		ExchangeDeployment:          test.StringPtr("test-value"),
		GoogleDeployment:            test.StringPtr("test-value"),
		GraphDeployment:             test.StringPtr("test-value"),
		ProcessAliasPrivateMeetings: true,
		ReplaceEmptySubject:         true,
		ReplaceSubjectType:          "template",
		ReplaceSubjectTemplate:      "test-value",
		UseWebex:                    true,
		WebexAPIDomain:              "test-value",
		WebexClientID:               test.StringPtr("test-value"),
		WebexClientSecret:           test.StringPtr("test-value"),
		WebexRedirectURI:            test.StringPtr("test-value"),
		WebexRefreshToken:           test.StringPtr("test-value"),
	}

	// Mock the GetMjxintegration API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_integration/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		mjxintegration := args.Get(2).(*config.MjxIntegration)
		*mjxintegration = *mockState
	}).Maybe()

	// Mock the UpdateMjxintegration API call
	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_integration/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.MjxIntegrationUpdateRequest)
		mjxintegration := args.Get(3).(*config.MjxIntegration)

		// Update mock state based on request
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.StartBuffer != nil {
			mockState.StartBuffer = *updateRequest.StartBuffer
		}
		if updateRequest.EndBuffer != nil {
			mockState.EndBuffer = *updateRequest.EndBuffer
		}
		if updateRequest.DisplayUpcomingMeetings != nil {
			mockState.DisplayUpcomingMeetings = *updateRequest.DisplayUpcomingMeetings
		}
		if updateRequest.EnableNonVideoMeetings != nil {
			mockState.EnableNonVideoMeetings = *updateRequest.EnableNonVideoMeetings
		}
		if updateRequest.EnablePrivateMeetings != nil {
			mockState.EnablePrivateMeetings = *updateRequest.EnablePrivateMeetings
		}
		if updateRequest.EPUseHTTPS != nil {
			mockState.EPUseHTTPS = *updateRequest.EPUseHTTPS
		}
		if updateRequest.EPVerifyCertificate != nil {
			mockState.EPVerifyCertificate = *updateRequest.EPVerifyCertificate
		}
		if updateRequest.ReplaceSubjectType != "" {
			mockState.ReplaceSubjectType = updateRequest.ReplaceSubjectType
		}
		if updateRequest.ExchangeDeployment != nil {
			mockState.ExchangeDeployment = updateRequest.ExchangeDeployment
		}
		if updateRequest.GoogleDeployment != nil {
			mockState.GoogleDeployment = updateRequest.GoogleDeployment
		}
		if updateRequest.GraphDeployment != nil {
			mockState.GraphDeployment = updateRequest.GraphDeployment
		}
		if updateRequest.ProcessAliasPrivateMeetings != nil {
			mockState.ProcessAliasPrivateMeetings = *updateRequest.ProcessAliasPrivateMeetings
		}
		if updateRequest.ReplaceEmptySubject != nil {
			mockState.ReplaceEmptySubject = *updateRequest.ReplaceEmptySubject
		}
		if updateRequest.ReplaceSubjectTemplate != "" {
			mockState.ReplaceSubjectTemplate = updateRequest.ReplaceSubjectTemplate
		}
		if updateRequest.UseWebex != nil {
			mockState.UseWebex = *updateRequest.UseWebex
		}
		if updateRequest.WebexAPIDomain != "" {
			mockState.WebexAPIDomain = updateRequest.WebexAPIDomain
		}
		if updateRequest.WebexClientID != nil {
			mockState.WebexClientID = updateRequest.WebexClientID
		}
		if updateRequest.WebexClientSecret != nil {
			mockState.WebexClientSecret = updateRequest.WebexClientSecret
		}
		if updateRequest.WebexRedirectURI != nil {
			mockState.WebexRedirectURI = updateRequest.WebexRedirectURI
		}
		if updateRequest.WebexRefreshToken != nil {
			mockState.WebexRefreshToken = updateRequest.WebexRefreshToken
		}
		if updateRequest.EPPassword != "" {
			mockState.EPPassword = updateRequest.EPPassword
		}

		// Return updated state
		*mjxintegration = *mockState
	}).Maybe()

	// Mock the DeleteMjxintegration API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/mjx_integration/123/"
	}), mock.Anything).Return(nil)

	testInfinityMjxIntegration(t, client)
}

func testInfinityMjxIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.mjx_integration-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.mjx_integration-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "name", "mjx_integration-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "description", "Test MjxIntegration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "enable_non_video_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "enable_private_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "ep_username", "mjx_integration-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "ep_use_https", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "ep_verify_certificate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "process_alias_private_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "replace_empty_subject", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "use_webex", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.mjx_integration-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.mjx_integration-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "name", "mjx_integration-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "description", "Updated Test MjxIntegration"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "enable_non_video_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "enable_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "ep_username", "mjx_integration-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "ep_use_https", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "ep_verify_certificate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "process_alias_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "replace_empty_subject", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.mjx_integration-test", "use_webex", "false"),
				),
			},
		},
	})
}
