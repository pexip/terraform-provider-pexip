package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/terraform-provider-pexip/internal/test"
	"os"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38"
	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/pexip/go-infinity-sdk/v38/types"
	"github.com/stretchr/testify/mock"
)

func TestInfinityNode(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateWorkerVM API call
	createResponse := &types.PostResponse{
		Body:        []byte("<xml><worker_vm><id>123</id></worker_vm></xml>"),
		ResourceURI: "/api/admin/configuration/v1/worker_vm/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/worker_vm/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Track state to return different values before and after update
	updated := false

	// Mock the GetWorkerVM API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/worker_vm/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		vm := args.Get(2).(*config.WorkerVM)
		if updated {
			*vm = config.WorkerVM{
				ID:                    123,
				Name:                  "test-node-1",
				Hostname:              "test-node-1",
				Address:               "192.168.1.100",
				Netmask:               "255.255.255.0",
				Domain:                "pexip.com",
				Gateway:               "192.168.1.1",
				Password:              "password123",
				NodeType:              "CONFERENCING",
				SystemLocation:        "Test Location",
				Transcoding:           true,
				VMCPUCount:            4,
				VMSystemMemory:        8192,
				MaintenanceMode:       false,
				MaintenanceModeReason: "",
				ResourceURI:           "/api/admin/configuration/v1/worker_vm/123/",
			}
		} else {
			*vm = config.WorkerVM{
				ID:                    123,
				Name:                  "test-node-1",
				Hostname:              "test-node-1",
				Address:               "192.168.1.100",
				Netmask:               "255.255.255.0",
				Domain:                "example.com",
				Gateway:               "192.168.1.1",
				Password:              "password123",
				NodeType:              "CONFERENCING",
				SystemLocation:        "Test Location",
				Transcoding:           true,
				VMCPUCount:            4,
				VMSystemMemory:        8192,
				MaintenanceMode:       false,
				MaintenanceModeReason: "",
				ResourceURI:           "/api/admin/configuration/v1/worker_vm/123/",
			}
		}
	}).Maybe() // Called multiple times for reads

	// Mock the UpdateWorkerVM API call
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/worker_vm/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		vm := args.Get(3).(*config.WorkerVM)
		*vm = config.WorkerVM{
			ID:                    123,
			Name:                  "test-node-1",
			Hostname:              "test-node-1",
			Address:               "192.168.1.100",
			Netmask:               "255.255.255.0",
			Domain:                "pexip.com",
			Gateway:               "192.168.1.1",
			Password:              "password123",
			NodeType:              "CONFERENCING",
			SystemLocation:        "Test Location",
			Transcoding:           true,
			VMCPUCount:            4,
			VMSystemMemory:        8192,
			MaintenanceMode:       false,
			MaintenanceModeReason: "",
			ResourceURI:           "/api/admin/configuration/v1/worker_vm/123/",
		}
	}).Maybe()

	// Mock the DeleteWorkerVM API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/worker_vm/123/"
	}), mock.Anything).Return(nil)

	testInfinityNode(t, client)
}

func testInfinityNode(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestData(t, "resource_infinity_node_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "config"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "name", "test-node-1"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "hostname", "test-node-1"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "node_type", "CONFERENCING"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "system_location", "Test Location"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "transcoding", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "domain", "example.com"),
				),
			},
			{
				Config: test.LoadTestData(t, "resource_infinity_node_basic_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "name"),
					resource.TestCheckResourceAttrSet("pexip_infinity_node.node", "config"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "name", "test-node-1"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "hostname", "test-node-1"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "node_type", "CONFERENCING"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "system_location", "Test Location"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "transcoding", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_node.node", "domain", "pexip.com"),
				),
			},
		},
	})
}
