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

func stringToMjxIntegrationRef(s *string) *config.MjxIntegrationResourceReference {
	if s == nil {
		return nil
	}
	return &config.MjxIntegrationResourceReference{ResourceURI: *s}
}

func TestInfinityMjxIntegration(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.MjxIntegration{}

	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/mjx_integration/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/mjx_integration/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.MjxIntegrationCreateRequest)
		*mockState = config.MjxIntegration{
			ID:                          123,
			ResourceURI:                 "/api/admin/configuration/v1/mjx_integration/123/",
			Name:                        createReq.Name,
			Description:                 createReq.Description,
			DisplayUpcomingMeetings:     createReq.DisplayUpcomingMeetings,
			EnableNonVideoMeetings:      createReq.EnableNonVideoMeetings,
			EnablePrivateMeetings:       createReq.EnablePrivateMeetings,
			EndBuffer:                   createReq.EndBuffer,
			StartBuffer:                 createReq.StartBuffer,
			EPUsername:                  createReq.EPUsername,
			EPUseHTTPS:                  createReq.EPUseHTTPS,
			EPVerifyCertificate:         createReq.EPVerifyCertificate,
			ExchangeDeployment:          stringToMjxIntegrationRef(createReq.ExchangeDeployment),
			GoogleDeployment:            stringToMjxIntegrationRef(createReq.GoogleDeployment),
			GraphDeployment:             stringToMjxIntegrationRef(createReq.GraphDeployment),
			ProcessAliasPrivateMeetings: createReq.ProcessAliasPrivateMeetings,
			ReplaceEmptySubject:         createReq.ReplaceEmptySubject,
			ReplaceSubjectType:          createReq.ReplaceSubjectType,
			ReplaceSubjectTemplate:      createReq.ReplaceSubjectTemplate,
			UseWebex:                    createReq.UseWebex,
			WebexAPIDomain:              createReq.WebexAPIDomain,
			WebexClientID:               createReq.WebexClientID,
			WebexOAuthState:             createReq.WebexOAuthState,
			WebexRedirectURI:            createReq.WebexRedirectURI,
		}
	})

	client.On("GetJSON", mock.Anything, "configuration/v1/mjx_integration/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		integration := args.Get(3).(*config.MjxIntegration)
		*integration = *mockState
	}).Maybe()

	client.On("PutJSON", mock.Anything, "configuration/v1/mjx_integration/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.MjxIntegrationUpdateRequest)
		mockState.Name = updateReq.Name
		mockState.Description = updateReq.Description
		mockState.DisplayUpcomingMeetings = updateReq.DisplayUpcomingMeetings
		mockState.EnableNonVideoMeetings = updateReq.EnableNonVideoMeetings
		mockState.EnablePrivateMeetings = updateReq.EnablePrivateMeetings
		mockState.EndBuffer = updateReq.EndBuffer
		mockState.StartBuffer = updateReq.StartBuffer
		mockState.EPUsername = updateReq.EPUsername
		mockState.EPUseHTTPS = updateReq.EPUseHTTPS
		mockState.EPVerifyCertificate = updateReq.EPVerifyCertificate
		mockState.ExchangeDeployment = stringToMjxIntegrationRef(updateReq.ExchangeDeployment)
		mockState.GoogleDeployment = stringToMjxIntegrationRef(updateReq.GoogleDeployment)
		mockState.GraphDeployment = stringToMjxIntegrationRef(updateReq.GraphDeployment)
		mockState.ProcessAliasPrivateMeetings = updateReq.ProcessAliasPrivateMeetings
		mockState.ReplaceEmptySubject = updateReq.ReplaceEmptySubject
		mockState.ReplaceSubjectType = updateReq.ReplaceSubjectType
		mockState.ReplaceSubjectTemplate = updateReq.ReplaceSubjectTemplate
		mockState.UseWebex = updateReq.UseWebex
		mockState.WebexAPIDomain = updateReq.WebexAPIDomain
		mockState.WebexClientID = updateReq.WebexClientID
		mockState.WebexOAuthState = updateReq.WebexOAuthState
		mockState.WebexRedirectURI = updateReq.WebexRedirectURI
	}).Maybe()

	client.On("DeleteJSON", mock.Anything, "configuration/v1/mjx_integration/123/", mock.Anything).Return(nil)

	testInfinityMjxIntegration(t, client)
}

func testInfinityMjxIntegration(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", "Test MJX integration description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "14"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_non_video_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "end_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "start_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_username", "ep-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_password", "ep-password-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_use_https", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_verify_certificate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "graph_deployment", "/api/admin/configuration/v1/mjx_graph_deployment/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "process_alias_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_empty_subject", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_type", "ALL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_template", "Meeting: {{ subject }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_api_domain", "custom.webexapis.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_client_id", "webex-client-id-full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_client_secret", "webex-secret-full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_redirect_uri", "https://pexip.example.com/admin/platform/mjxintegration/oauth_redirect/"),
				),
			},
			// Step 2: Update with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "7"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_non_video_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_private_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "end_buffer", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "start_buffer", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_use_https", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_verify_certificate", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "graph_deployment", "/api/admin/configuration/v1/mjx_graph_deployment/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "process_alias_private_meetings", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_empty_subject", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_type", "PRIVATE"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_template", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_api_domain", "webexapis.com"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_mjx_integration_min"),
				Destroy: true,
			},
			// Step 4: Recreate with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration min"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "7"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "graph_deployment", "/api/admin/configuration/v1/mjx_graph_deployment/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "false"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_mjx_integration_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_mjx_integration.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "name", "tf-test mjx-integration full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "description", "Test MJX integration description"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "display_upcoming_meetings", "14"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_non_video_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "enable_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "end_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "start_buffer", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_username", "ep-user@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_password", "ep-password-test"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_use_https", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "ep_verify_certificate", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "graph_deployment", "/api/admin/configuration/v1/mjx_graph_deployment/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "process_alias_private_meetings", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_empty_subject", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_type", "ALL"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "replace_subject_template", "Meeting: {{ subject }}"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "use_webex", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_api_domain", "custom.webexapis.com"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_client_id", "webex-client-id-full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_client_secret", "webex-secret-full"),
					resource.TestCheckResourceAttr("pexip_infinity_mjx_integration.test", "webex_redirect_uri", "https://pexip.example.com/admin/platform/mjxintegration/oauth_redirect/"),
				),
			},
		},
	})
}
