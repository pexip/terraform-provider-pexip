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

func TestInfinityScheduledScaling(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Shared state for mocking
	mockState := &config.ScheduledScaling{
		ID:                 123,
		ResourceURI:        "/api/admin/configuration/v1/scheduled_scaling/123/",
		PolicyName:         "scheduled_scaling-test",
		PolicyType:         "TeamsConnectorScaling",
		ResourceIdentifier: "tf-test teams-proxy scheduled scaling",
		Enabled:            true,
		LocalTimezone:      "UTC",
		StartDate:          "2024-01-01",
		TimeFrom:           "09:00:00",
		TimeTo:             "17:00:00",
		InstancesToAdd:     0,
		MinutesInAdvance:   20,
		Mon:                true,
		Tue:                false,
		Wed:                false,
		Thu:                false,
		Fri:                false,
		Sat:                false,
		Sun:                false,
	}

	// Shared Azure Tenant state
	azureTenantState := &config.AzureTenant{
		ID:          456,
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/456/",
		Name:        "tf-test azure-tenant-teams-proxy scheduled scaling",
		Description: "Test Azure Tenant for Scheduled Scaling",
		TenantID:    "44444444-4444-4444-4444-444444444445",
	}

	// Mock Azure Tenant create
	azureTenantCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/456/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/azure_tenant/", mock.Anything, mock.Anything).Return(azureTenantCreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.AzureTenantCreateRequest)
		azureTenantState.Name = createReq.Name
		azureTenantState.Description = createReq.Description
		azureTenantState.TenantID = createReq.TenantID
	}).Maybe()

	// Mock Azure Tenant read
	client.On("GetJSON", mock.Anything, "configuration/v1/azure_tenant/456/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		azureTenant := args.Get(3).(*config.AzureTenant)
		*azureTenant = *azureTenantState
	}).Maybe()

	// Mock Azure Tenant update
	client.On("PutJSON", mock.Anything, "configuration/v1/azure_tenant/456/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.AzureTenantUpdateRequest)
		azureTenant := args.Get(3).(*config.AzureTenant)
		azureTenantState.Name = updateReq.Name
		azureTenantState.Description = updateReq.Description
		azureTenantState.TenantID = updateReq.TenantID
		*azureTenant = *azureTenantState
	}).Maybe()

	// Mock Azure Tenant delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/azure_tenant/456/", mock.Anything).Return(nil).Maybe()

	// Shared Teams Proxy state
	teamsProxyState := &config.TeamsProxy{
		ID:                   789,
		ResourceURI:          "/api/admin/configuration/v1/teams_proxy/789/",
		Name:                 "tf-test teams-proxy scheduled scaling",
		Address:              "teams-proxy-min.pexvclab.com",
		Port:                 443,
		AzureTenant:          "/api/admin/configuration/v1/azure_tenant/456/",
		MinNumberOfInstances: 1,
		NotificationsEnabled: false,
	}

	// Mock Teams Proxy create
	teamsProxyCreateResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/teams_proxy/789/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/teams_proxy/", mock.Anything, mock.Anything).Return(teamsProxyCreateResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.TeamsProxyCreateRequest)
		teamsProxyState.Name = createReq.Name
		teamsProxyState.Address = createReq.Address
		teamsProxyState.Port = createReq.Port
		teamsProxyState.AzureTenant = createReq.AzureTenant
		teamsProxyState.MinNumberOfInstances = createReq.MinNumberOfInstances
	}).Maybe()

	// Mock Teams Proxy read
	client.On("GetJSON", mock.Anything, "configuration/v1/teams_proxy/789/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		teamsProxy := args.Get(3).(*config.TeamsProxy)
		*teamsProxy = *teamsProxyState
	}).Maybe()

	// Mock Teams Proxy update
	client.On("PutJSON", mock.Anything, "configuration/v1/teams_proxy/789/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.TeamsProxyUpdateRequest)
		teamsProxy := args.Get(3).(*config.TeamsProxy)
		teamsProxyState.Name = updateReq.Name
		teamsProxyState.Address = updateReq.Address
		teamsProxyState.Port = updateReq.Port
		teamsProxyState.AzureTenant = updateReq.AzureTenant
		teamsProxyState.MinNumberOfInstances = updateReq.MinNumberOfInstances
		*teamsProxy = *teamsProxyState
	}).Maybe()

	// Mock Teams Proxy delete
	client.On("DeleteJSON", mock.Anything, "configuration/v1/teams_proxy/789/", mock.Anything).Return(nil).Maybe()

	// Mock the CreateScheduledscaling API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/scheduled_scaling/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/scheduled_scaling/", mock.Anything, mock.Anything).Return(createResponse, nil).Run(func(args mock.Arguments) {
		createReq := args.Get(2).(*config.ScheduledScalingCreateRequest)
		// Update mockState with values from create request
		mockState.PolicyName = createReq.PolicyName
		mockState.PolicyType = createReq.PolicyType
		mockState.ResourceIdentifier = createReq.ResourceIdentifier
		mockState.Enabled = createReq.Enabled
		mockState.LocalTimezone = createReq.LocalTimezone
		mockState.StartDate = createReq.StartDate
		mockState.TimeFrom = createReq.TimeFrom
		mockState.TimeTo = createReq.TimeTo
		mockState.InstancesToAdd = createReq.InstancesToAdd
		mockState.MinutesInAdvance = createReq.MinutesInAdvance
		mockState.Mon = createReq.Mon
		mockState.Tue = createReq.Tue
		mockState.Wed = createReq.Wed
		mockState.Thu = createReq.Thu
		mockState.Fri = createReq.Fri
		mockState.Sat = createReq.Sat
		mockState.Sun = createReq.Sun
	})

	// Mock the GetScheduledscaling API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/scheduled_scaling/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		scheduled_scaling := args.Get(3).(*config.ScheduledScaling)
		*scheduled_scaling = *mockState
	}).Maybe()

	// Mock the UpdateScheduledscaling API call
	client.On("PutJSON", mock.Anything, "configuration/v1/scheduled_scaling/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.ScheduledScalingUpdateRequest)
		scheduled_scaling := args.Get(3).(*config.ScheduledScaling)

		// Update mock state based on request - UpdateRequest fields are not pointers per SDK
		mockState.PolicyName = updateReq.PolicyName
		mockState.PolicyType = updateReq.PolicyType
		mockState.ResourceIdentifier = updateReq.ResourceIdentifier
		mockState.LocalTimezone = updateReq.LocalTimezone
		mockState.StartDate = updateReq.StartDate
		mockState.TimeFrom = updateReq.TimeFrom
		mockState.TimeTo = updateReq.TimeTo

		if updateReq.Enabled != nil {
			mockState.Enabled = *updateReq.Enabled
		}
		if updateReq.InstancesToAdd != nil {
			mockState.InstancesToAdd = *updateReq.InstancesToAdd
		}
		if updateReq.MinutesInAdvance != nil {
			mockState.MinutesInAdvance = *updateReq.MinutesInAdvance
		}
		if updateReq.Mon != nil {
			mockState.Mon = *updateReq.Mon
		}
		if updateReq.Tue != nil {
			mockState.Tue = *updateReq.Tue
		}
		if updateReq.Wed != nil {
			mockState.Wed = *updateReq.Wed
		}
		if updateReq.Thu != nil {
			mockState.Thu = *updateReq.Thu
		}
		if updateReq.Fri != nil {
			mockState.Fri = *updateReq.Fri
		}
		if updateReq.Sat != nil {
			mockState.Sat = *updateReq.Sat
		}
		if updateReq.Sun != nil {
			mockState.Sun = *updateReq.Sun
		}

		// Return updated state
		*scheduled_scaling = *mockState
	}).Maybe()

	// Mock the DeleteScheduledscaling API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/scheduled_scaling/123/"
	}), mock.Anything).Return(nil)

	testInfinityScheduledScaling(t, client)
}

func testInfinityScheduledScaling(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				// Step 1: Create with full config
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_full"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_name", "tf-test scheduled scaling full"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_type", "TeamsConnectorScaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "resource_identifier", "tf-test teams-proxy scheduled scaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "local_timezone", "America/New_York"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "start_date", "2024-06-15"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_from", "08:30:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_to", "18:30:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "instances_to_add", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "minutes_in_advance", "30"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "mon", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "tue", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "wed", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "thu", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "fri", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sat", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sun", "true"),
				),
			},
			{
				// Step 2: Update to min config
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_name", "tf-test scheduled scaling min"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_type", "TeamsConnectorScaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "resource_identifier", "tf-test teams-proxy scheduled scaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "local_timezone", "UTC"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "start_date", "2024-01-01"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_from", "09:00:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_to", "17:00:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "instances_to_add", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "minutes_in_advance", "20"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "mon", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "tue", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "wed", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "thu", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "fri", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sat", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sun", "false"),
				),
			},
			{
				// Step 3: Destroy and recreate with minimal config
				Config:  test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_min"),
				Destroy: true,
			},
			{
				// Step 4: Recreate with minimal config (after destroy)
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_min"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_name", "tf-test scheduled scaling min"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_type", "TeamsConnectorScaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "resource_identifier", "tf-test teams-proxy scheduled scaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "enabled", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "local_timezone", "UTC"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "start_date", "2024-01-01"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_from", "09:00:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_to", "17:00:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "instances_to_add", "0"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "minutes_in_advance", "20"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "mon", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "tue", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "wed", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "thu", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "fri", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sat", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sun", "false"),
				),
			},
			{
				// Step 5: Update to full config
				Config: test.LoadTestFolder(t, "resource_infinity_scheduled_scaling_full"),
				Check: resource.ComposeTestCheckFunc(
					// IDs and required fields
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_scheduled_scaling.test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_name", "tf-test scheduled scaling full"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "policy_type", "TeamsConnectorScaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "resource_identifier", "tf-test teams-proxy scheduled scaling"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "enabled", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "local_timezone", "America/New_York"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "start_date", "2024-06-15"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_from", "08:30:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "time_to", "18:30:00"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "instances_to_add", "5"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "minutes_in_advance", "30"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "mon", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "tue", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "wed", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "thu", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "fri", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sat", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_scheduled_scaling.test", "sun", "true"),
				),
			},
		},
	})
}
