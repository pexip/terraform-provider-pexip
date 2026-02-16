/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

	// Mock azure_tenant creation
	azureTenantCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/456/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/azure_tenant/", mock.Anything, mock.Anything).Return(azureTenantCreateResponse, nil)

	// Shared state for azure_tenant mock
	azureTenantMockState := &config.AzureTenant{
		ID:          456,
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/456/",
		Name:        "tf-test-azure-tenant-teams-proxy-full",
		Description: "Test Azure Tenant for Teams Proxy",
		TenantID:    "12345678-1234-1234-1234-123456789012",
	}

	// Mock azure_tenant read
	client.On("GetJSON", mock.Anything, "configuration/v1/azure_tenant/456/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		azureTenant := args.Get(3).(*config.AzureTenant)
		*azureTenant = *azureTenantMockState
	}).Maybe()

	// Mock azure_tenant update
	client.On("PutJSON", mock.Anything, "configuration/v1/azure_tenant/456/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.AzureTenantUpdateRequest)
		azureTenant := args.Get(3).(*config.AzureTenant)

		azureTenantMockState.Name = updateRequest.Name
		azureTenantMockState.Description = updateRequest.Description
		azureTenantMockState.TenantID = updateRequest.TenantID

		*azureTenant = *azureTenantMockState
	}).Maybe()

	// Mock azure_tenant delete
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/azure_tenant/456/"
	}), mock.Anything).Return(nil)

	// Mock teams_proxy creation
	teamsProxyCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/teams_proxy/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/teams_proxy/", mock.Anything, mock.Anything).Return(teamsProxyCreateResponse, nil)

	// Shared state for teams_proxy mock - dynamically updated by create/update operations
	teamsProxyMockState := &config.TeamsProxy{
		ID:                   123,
		ResourceURI:          "/api/admin/configuration/v1/teams_proxy/123/",
		Name:                 "tf-test-teams-proxy-full",
		Description:          "Test Teams Proxy Full Configuration",
		Address:              "teams-proxy-full.pexvclab.com",
		Port:                 8443,
		AzureTenant:          "/api/admin/configuration/v1/azure_tenant/456/",
		MinNumberOfInstances: 3,
		NotificationsEnabled: true,
		NotificationsQueue:   test.StringPtr("Endpoint=sb://examplevmss.servicebus.windows.net/;SharedAccessKeyName=standard_access_policy;SharedAccessKey=testkey123="),
		EventhubID:           test.StringPtr("akf445PWGCMdBEUUiHBXicMn4DQ="), // Computed by API from notifications_queue
	}

	// Mock teams_proxy read
	client.On("GetJSON", mock.Anything, "configuration/v1/teams_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		teamsProxy := args.Get(3).(*config.TeamsProxy)
		*teamsProxy = *teamsProxyMockState
	}).Maybe()

	// Mock teams_proxy update
	client.On("PutJSON", mock.Anything, "configuration/v1/teams_proxy/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.TeamsProxyUpdateRequest)
		teamsProxy := args.Get(3).(*config.TeamsProxy)

		// Update mock state with all fields from request
		teamsProxyMockState.Name = updateRequest.Name
		teamsProxyMockState.Description = updateRequest.Description
		teamsProxyMockState.Address = updateRequest.Address
		teamsProxyMockState.Port = updateRequest.Port
		teamsProxyMockState.AzureTenant = updateRequest.AzureTenant
		teamsProxyMockState.MinNumberOfInstances = updateRequest.MinNumberOfInstances
		teamsProxyMockState.NotificationsEnabled = updateRequest.NotificationsEnabled
		teamsProxyMockState.NotificationsQueue = updateRequest.NotificationsQueue

		// Simulate API behavior: eventhub_id is computed from notifications_queue and persists
		if updateRequest.NotificationsQueue != nil && *updateRequest.NotificationsQueue != "" {
			teamsProxyMockState.EventhubID = test.StringPtr("akf445PWGCMdBEUUiHBXicMn4DQ=")
		}
		// Note: eventhub_id persists even when notifications_queue is cleared (API behavior)

		// Return updated state
		*teamsProxy = *teamsProxyMockState
	}).Maybe()

	// Mock teams_proxy delete
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
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "tf-test-teams-proxy-full"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", "Test Teams Proxy Full Configuration"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "teams-proxy-full.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "8443"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "3"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", "Endpoint=sb://examplevmss.servicebus.windows.net/;SharedAccessKeyName=standard_access_policy;SharedAccessKey=testkey123="),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id"), // Computed by API
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "tf-test-teams-proxy-min"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "teams-proxy-min.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "443"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id"), // Persists from previous config
				),
			},
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_teams_proxy_min"),
				Destroy: true,
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "tf-test-teams-proxy-min"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "teams-proxy-min.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "443"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", ""),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id"), // Persists from previous config
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_teams_proxy_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "name", "tf-test-teams-proxy-full"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "description", "Test Teams Proxy Full Configuration"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "address", "teams-proxy-full.pexvclab.com"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "port", "8443"),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "azure_tenant"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "min_number_of_instances", "3"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_teams_proxy.teams-proxy-test", "notifications_queue", "Endpoint=sb://examplevmss.servicebus.windows.net/;SharedAccessKeyName=standard_access_policy;SharedAccessKey=testkey123="),
					resource.TestCheckResourceAttrSet("pexip_infinity_teams_proxy.teams-proxy-test", "eventhub_id"), // Computed by API
				),
			},
		},
	})
}
