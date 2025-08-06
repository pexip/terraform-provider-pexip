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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityWorkerVM(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateWorkervm API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/worker_vm/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/worker_vm/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.WorkerVM{
		ID:                    123,
		ResourceURI:           "/api/admin/configuration/v1/worker_vm/123/",
		Name:                  "worker_vm-test",
		Hostname:              "worker_vm-test",
		Domain:                "test-value",
		Address:               "192.168.1.10",
		Netmask:               "255.255.255.0",
		Gateway:               "192.168.1.1",
		IPv6Address:           test.StringPtr("2001:db8::1"),
		IPv6Gateway:           test.StringPtr("2001:db8::fe"),
		NodeType:              "CONFERENCING",
		Transcoding:           true,
		Password:              "test-value",
		MaintenanceMode:       true,
		MaintenanceModeReason: "test-value",
		SystemLocation:        "test-value",
		VMCPUCount:            4,
		VMSystemMemory:        4096,
	}

	// Mock the GetWorkervm API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/worker_vm/123/", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		worker_vm := args.Get(2).(*config.WorkerVM)
		*worker_vm = *mockState
	}).Maybe()

	// Mock the UpdateWorkervm API call
	client.On("PutJSON", mock.Anything, "configuration/v1/worker_vm/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.WorkerVMUpdateRequest)
		worker_vm := args.Get(3).(*config.WorkerVM)

		// Update mock state based on request
		if updateRequest.Name != "" {
			mockState.Name = updateRequest.Name
		}
		if updateRequest.Hostname != "" {
			mockState.Hostname = updateRequest.Hostname
		}
		if updateRequest.Domain != "" {
			mockState.Domain = updateRequest.Domain
		}
		if updateRequest.Address != "" {
			mockState.Address = updateRequest.Address
		}
		if updateRequest.Netmask != "" {
			mockState.Netmask = updateRequest.Netmask
		}
		if updateRequest.Gateway != "" {
			mockState.Gateway = updateRequest.Gateway
		}
		if updateRequest.IPv6Address != nil {
			mockState.IPv6Address = updateRequest.IPv6Address
		}
		if updateRequest.IPv6Gateway != nil {
			mockState.IPv6Gateway = updateRequest.IPv6Gateway
		}
		if updateRequest.NodeType != "" {
			mockState.NodeType = updateRequest.NodeType
		}
		mockState.Transcoding = updateRequest.Transcoding
		if updateRequest.Password != "" {
			mockState.Password = updateRequest.Password
		}
		mockState.MaintenanceMode = updateRequest.MaintenanceMode
		if updateRequest.MaintenanceModeReason != "" {
			mockState.MaintenanceModeReason = updateRequest.MaintenanceModeReason
		}
		if updateRequest.SystemLocation != "" {
			mockState.SystemLocation = updateRequest.SystemLocation
		}

		// Return updated state
		*worker_vm = *mockState
	}).Maybe()

	// Mock the DeleteWorkervm API call
	client.On("DeleteJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/worker_vm/123/"
	}), mock.Anything).Return(nil)

	testInfinityWorkerVM(t, client)
}

func testInfinityWorkerVM(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker_vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker_vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "name", "worker_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "hostname", "worker_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "domain", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "transcoding", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "maintenance_mode", "true"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker_vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker_vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "name", "worker_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "hostname", "worker_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "domain", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "address", "192.168.1.20"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "transcoding", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker_vm-test", "maintenance_mode", "false"),
				),
			},
		},
	})
}
