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

func TestInfinityAzureTenant(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateAzureTenant API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/azure_tenant/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.AzureTenant{
		ID:          123,
		ResourceURI: "/api/admin/configuration/v1/azure_tenant/123/",
		Name:        "azure_tenant-test",
		Description: "Test AzureTenant",
		TenantID:    "12345678-1234-1234-1234-123456789012",
	}

	// Mock the GetAzureTenant API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/azure_tenant/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		azureTenant := args.Get(3).(*config.AzureTenant)
		*azureTenant = *mockState
	}).Maybe()

	// Mock the UpdateAzureTenant API call
	client.On("PutJSON", mock.Anything, "configuration/v1/azure_tenant/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.AzureTenantUpdateRequest)
		azureTenant := args.Get(3).(*config.AzureTenant)

		// Update mock state
		mockState.Name = updateRequest.Name
		if updateRequest.Description != "" {
			mockState.Description = updateRequest.Description
		}
		if updateRequest.TenantID != "" {
			mockState.TenantID = updateRequest.TenantID
		}

		// Return updated state
		*azureTenant = *mockState
	}).Maybe()

	// Mock the DeleteAzureTenant API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/azure_tenant/123/"
	}), mock.Anything).Return(nil)

	testInfinityAzureTenant(t, client)
}

func testInfinityAzureTenant(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_azure_tenant_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_azure_tenant.azure_tenant-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_azure_tenant.azure_tenant-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_azure_tenant.azure_tenant-test", "name", "azure_tenant-test"),
					resource.TestCheckResourceAttr("pexip_infinity_azure_tenant.azure_tenant-test", "description", "Test AzureTenant"),
					resource.TestCheckResourceAttr("pexip_infinity_azure_tenant.azure_tenant-test", "tenant_id", "12345678-1234-1234-1234-123456789012"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_azure_tenant_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_azure_tenant.azure_tenant-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_azure_tenant.azure_tenant-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_azure_tenant.azure_tenant-test", "name", "azure_tenant-test"),
					resource.TestCheckResourceAttr("pexip_infinity_azure_tenant.azure_tenant-test", "description", "Updated Test AzureTenant"),
					resource.TestCheckResourceAttr("pexip_infinity_azure_tenant.azure_tenant-test", "tenant_id", "87654321-4321-4321-4321-210987654321"),
				),
			},
		},
	})
}
