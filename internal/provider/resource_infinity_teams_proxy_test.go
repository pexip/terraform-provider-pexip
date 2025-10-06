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

	// Track state to return different values before and after update
	updated := false

	// config for create request with minimal fields set
	// include all required fields and fields with non-null defaults
	createMockState := &config.TeamsProxy{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/teams_proxy/123/",
		Name:                 "test-teams-proxy",
		Description:          "",
		Address:              "test-teams-proxy.dev.pexip.network",
		Port:                 443,
		AzureTenant:          "/api/admin/configuration/v1/azure_tenant/1/",
		MinNumberOfInstances: 1,
		NotificationsEnabled: false,
	}

	// config for update request with all fields set and updated
	updateMockState := &config.TeamsProxy{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/teams_proxy/123/",
		Name:                 "test-teams-proxy-updated",
		Description:          "Test Teams Proxy Updated",
		Address:              "updated-test-teams-proxy.dev.pexip.network",
		Port:                 8443,
		AzureTenant:          "/api/admin/configuration/v1/azure_tenant/1/",
		MinNumberOfInstances: 0,
		NotificationsEnabled: true,
		NotificationsQueue:   test.StringPtr("updated-test-notifications-queue"),
		EventhubID:           test.StringPtr("updated-test-eventhub-id"),
	}

	// Mock the GetTeamsProxy API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/teams_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		teamsProxy := args.Get(3).(*config.TeamsProxy)
		if updated {
			*teamsProxy = *updateMockState
		} else {
			*teamsProxy = *createMockState
		}
	}).Maybe()

	// Mock the UpdateTeamsProxy API call
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/teams_proxy/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		teamsProxy := args.Get(3).(*config.TeamsProxy)
		*teamsProxy = *updateMockState
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
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "test-teams-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "443"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant", "/api/admin/configuration/v1/azure_tenant/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "false"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "test-teams-proxy-updated"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", "Test Teams Proxy Updated"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "updated-test-teams-proxy.dev.pexip.network"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "8443"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant", "/api/admin/configuration/v1/azure_tenant/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id", "updated-test-eventhub-id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", "updated-test-notifications-queue"),
				),
			},
		},
	})
}
