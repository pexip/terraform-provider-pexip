/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityTeamsProxy(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateTeamsProxy API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/teams_proxy/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/teams_proxy/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	eventhubID := "test-eventhub-id"
	notificationsQueue := "test-notifications-queue"
	mockState := &config.TeamsProxy{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/teams_proxy/123/",
		Name:                 "test-teams-proxy",
		Description:          "Test Teams Proxy",
		Address:              "test-teams-proxy.dev.pexip.network",
		Port:                 8080,
		AzureTenant:          "test-azure-tenant",
		EventhubID:           &eventhubID,
		MinNumberOfInstances: 2,
		NotificationsEnabled: false,
		NotificationsQueue:   &notificationsQueue,
	}

	// Mock the GetTeamsProxy API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/teams_proxy/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		teamsProxy := args.Get(2).(*config.TeamsProxy)
		*teamsProxy = *mockState
	}).Maybe()

	// Mock the UpdateTeamsProxy API call
	client.On("PutJSON", mock.Anything, "configuration/v1/teams_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.TeamsProxyUpdateRequest)
		teamsProxy := args.Get(3).(*config.TeamsProxy)

		// Update mock state
		mockState.Name = updateRequest.Name
		mockState.Description = updateRequest.Description
		mockState.Address = updateRequest.Address
		mockState.AzureTenant = updateRequest.AzureTenant
		if updateRequest.Port != nil {
			mockState.Port = *updateRequest.Port
		}
		if updateRequest.EventhubID != nil {
			mockState.EventhubID = updateRequest.EventhubID
		}
		if updateRequest.MinNumberOfInstances != nil {
			mockState.MinNumberOfInstances = *updateRequest.MinNumberOfInstances
		}
		if updateRequest.NotificationsEnabled != nil {
			mockState.NotificationsEnabled = *updateRequest.NotificationsEnabled
		}
		if updateRequest.NotificationsQueue != nil {
			mockState.NotificationsQueue = updateRequest.NotificationsQueue
		}

		// Return updated state
		*teamsProxy = *mockState
	}).Maybe()

	// Mock the DeleteTeamsProxy API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/teams_proxy/123/"
	}), mock.Anything).Return(nil)

	testInfinityTeamsProxy(t, client)
}

func testInfinityTeamsProxy(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "test-teams-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", "Test Teams Proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "test-teams-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "8080"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant", "test-azure-tenant"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id", "test-eventhub-id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "2"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", "test-notifications-queue"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "test-teams-proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", "Test Teams Proxy"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "test-teams-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "8081"), // Updated port
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant", "test-azure-tenant"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id", "test-eventhub-id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "1"), // Updated min_number_of_instances
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", "test-notifications-queue"),
				),
			},
		},
	})
}
