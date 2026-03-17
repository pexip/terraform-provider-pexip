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

func TestInfinitySnmpNetworkManagementSystem(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateSnmpnetworkmanagementsystem API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/snmp_network_management_system/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/snmp_network_management_system/", mock.Anything, mock.Anything).Return(createResponse, nil).Maybe()

	// Shared state for mocking - starts with full config
	mockState := &config.SnmpNetworkManagementSystem{
		ID:                123,
		ResourceURI:       "/api/admin/configuration/v1/snmp_network_management_system/123/",
		Name:              "tf-test-snmp-nms",
		Description:       "tf-test SNMP NMS Description",
		Address:           "192.168.1.100",
		Port:              162,
		SnmpTrapCommunity: "tf-test-comm",
	}

	// Mock the GetSnmpnetworkmanagementsystem API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/snmp_network_management_system/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		snmp_network_management_system := args.Get(3).(*config.SnmpNetworkManagementSystem)
		*snmp_network_management_system = *mockState
	}).Maybe()

	// Mock the UpdateSnmpnetworkmanagementsystem API call
	client.On("PutJSON", mock.Anything, "configuration/v1/snmp_network_management_system/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SnmpNetworkManagementSystemUpdateRequest)
		snmp_network_management_system := args.Get(3).(*config.SnmpNetworkManagementSystem)

		// Update mock state based on request
		// Description and SnmpTrapCommunity always sent now (without omitempty)
		mockState.Description = updateReq.Description
		mockState.SnmpTrapCommunity = updateReq.SnmpTrapCommunity

		if updateReq.Address != "" {
			mockState.Address = updateReq.Address
		}
		if updateReq.Port != nil {
			mockState.Port = *updateReq.Port
		}

		// Return updated state
		*snmp_network_management_system = *mockState
	}).Maybe()

	// Mock the DeleteSnmpnetworkmanagementsystem API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/snmp_network_management_system/123/"
	}), mock.Anything).Return(nil).Maybe()

	testInfinitySnmpNetworkManagementSystem(t, client)
}

func testInfinitySnmpNetworkManagementSystem(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Create with full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "name", "tf-test-snmp-nms"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "description", "tf-test SNMP NMS Description"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "port", "162"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "snmp_trap_community", "tf-test-comm"),
				),
			},
			// Step 2: Update to min config (clearing optional fields)
			{
				Config: test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "name", "tf-test-snmp-nms"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "port", "161"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "snmp_trap_community", "public"),
				),
			},
			// Step 3: Destroy
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_min"),
				Destroy: true,
			},
			// Step 4: Create with min config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "name", "tf-test-snmp-nms"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "port", "161"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "snmp_trap_community", "public"),
				),
			},
			// Step 5: Update to full config
			{
				Config: test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "name", "tf-test-snmp-nms"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "description", "tf-test SNMP NMS Description"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "port", "162"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.tf-test-snmp-nms", "snmp_trap_community", "tf-test-comm"),
				),
			},
		},
	})
}
