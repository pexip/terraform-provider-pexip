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
	client.On("PostWithResponse", mock.Anything, "configuration/v1/snmp_network_management_system/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.SnmpNetworkManagementSystem{
		ID:                123,
		ResourceURI:       "/api/admin/configuration/v1/snmp_network_management_system/123/",
		Name:              "snmp_network_management_system-test",
		Description:       "Test SnmpNetworkManagementSystem",
		Address:           "192.168.1.100",
		Port:              162,
		SnmpTrapCommunity: "test-value",
	}

	// Mock the GetSnmpnetworkmanagementsystem API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/snmp_network_management_system/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		snmp_network_management_system := args.Get(2).(*config.SnmpNetworkManagementSystem)
		*snmp_network_management_system = *mockState
	}).Maybe()

	// Mock the UpdateSnmpnetworkmanagementsystem API call
	client.On("PutJSON", mock.Anything, "configuration/v1/snmp_network_management_system/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateReq := args.Get(2).(*config.SnmpNetworkManagementSystemUpdateRequest)
		snmp_network_management_system := args.Get(3).(*config.SnmpNetworkManagementSystem)

		// Update mock state based on request
		if updateReq.Description != "" {
			mockState.Description = updateReq.Description
		}
		if updateReq.Address != "" {
			mockState.Address = updateReq.Address
		}
		if updateReq.Port != nil {
			mockState.Port = *updateReq.Port
		}
		if updateReq.SnmpTrapCommunity != "" {
			mockState.SnmpTrapCommunity = updateReq.SnmpTrapCommunity
		}

		// Return updated state
		*snmp_network_management_system = *mockState
	}).Maybe()

	// Mock the DeleteSnmpnetworkmanagementsystem API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/snmp_network_management_system/123/"
	}), mock.Anything).Return(nil)

	testInfinitySnmpNetworkManagementSystem(t, client)
}

func testInfinitySnmpNetworkManagementSystem(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "name", "snmp_network_management_system-test"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "description", "Test SnmpNetworkManagementSystem"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_snmp_network_management_system_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "name", "snmp_network_management_system-test"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "description", "Updated Test SnmpNetworkManagementSystem"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "address", "192.168.1.200"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "port", "161"),
					resource.TestCheckResourceAttr("pexip_infinity_snmp_network_management_system.snmp_network_management_system-test", "snmp_trap_community", "updated-value"),
				),
			},
		},
	})
}
