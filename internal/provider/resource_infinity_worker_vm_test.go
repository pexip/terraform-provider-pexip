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

	// Track state to return different values before and after update
	updated := false

	// config for update that clears all fields
	// include all required fields and fields with non-null defaults
	clearAllFieldsMockState := &config.WorkerVM{
		ID:                         123,
		ResourceURI:                "/api/admin/configuration/v1/worker_vm/123/",
		Name:                       "worker-vm-test",
		Hostname:                   "worker-vm-test",
		Domain:                     "updated-value",
		Address:                    "192.168.1.10",
		Netmask:                    "255.255.255.0",
		Gateway:                    "192.168.1.1",
		SystemLocation:             "/api/admin/configuration/v1/system_location/200/",
		Description:                "",
		NodeType:                   "CONFERENCING",
		DeploymentType:             "MANUAL-PROVISION-ONLY",
		Transcoding:                true,
		MaintenanceMode:            false,
		MaintenanceModeReason:      "",
		MediaPriorityWeight:        test.IntPtr(0),
		VMCPUCount:                 4,
		VMSystemMemory:             4096,
		SSHAuthorizedKeysUseCloud:  true,
		SNMPAuthenticationPassword: "",
		SNMPCommunity:              "public",
		SNMPMode:                   "DISABLED",
		SNMPPrivacyPassword:        "",
		SNMPSystemContact:          "admin@domain.com",
		SNMPSystemLocation:         "Virtual machine",
		SNMPUsername:               "",
		EnableSSH:                  "GLOBAL",
		EnableDistributedDatabase:  true,
		CloudBursting:              false,
		ServiceManager:             true,
		ServicePolicy:              true,
		Signalling:                 true,
		Managed:                    false,
	}

	// Mock the GetWorkervm API call for Read operations
	client.On("GetJSON", mock.Anything, "configuration/v1/worker_vm/123/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		workerVM := args.Get(3).(*config.WorkerVM)
		if updated {
			*workerVM = *clearAllFieldsMockState
		} else {
			*workerVM = config.WorkerVM{
				ID:                    123,
				ResourceURI:           "/api/admin/configuration/v1/worker_vm/123/",
				Name:                  "worker-vm-test",
				Hostname:              "worker-vm-test",
				Domain:                "test-value",
				Address:               "192.168.1.10",
				Netmask:               "255.255.255.0",
				Gateway:               "192.168.1.1",
				SystemLocation:        "/api/admin/configuration/v1/system_location/1/",
				TLSCertificate:        test.StringPtr("/api/admin/configuration/v1/tls_certificate/2/"),
				Description:           "initial description",
				IPv6Address:           test.StringPtr("2001:db8::1"),
				IPv6Gateway:           test.StringPtr("2001:db8::fe"),
				NodeType:              "CONFERENCING",
				DeploymentType:        "MANUAL-PROVISION-ONLY",
				Transcoding:           true,
				Password:              "test-value",
				MaintenanceMode:       true,
				MaintenanceModeReason: "test-value",
				VMCPUCount:            4,
				VMSystemMemory:        4096,
				SecondaryAddress:      test.StringPtr("172.16.0.10"),
				SecondaryNetmask:      test.StringPtr("255.255.255.0"),
				MediaPriorityWeight:   test.IntPtr(10),
				SSHAuthorizedKeys: []string{
					"/api/admin/configuration/v1/ssh_authorized_key/1/",
				},
				SSHAuthorizedKeysUseCloud: true,
				StaticNATAddress:          test.StringPtr("203.0.113.2"),
				StaticRoutes: []string{
					"/api/admin/configuration/v1/static_route/1/",
				},
				SNMPAuthenticationPassword: "auth-password1",
				SNMPCommunity:              "public1",
				SNMPMode:                   "STANDARD",
				SNMPPrivacyPassword:        "privacy-password1",
				SNMPSystemContact:          "snmpcontact1@domain.com",
				SNMPSystemLocation:         "test-value",
				SNMPUsername:               "snmp-user1",
				EnableSSH:                  "ON",
				EnableDistributedDatabase:  false,
				ServiceManager:             true,
				ServicePolicy:              true,
				Signalling:                 true,
				Managed:                    false,
			}
		}
	}).Maybe()

	// Mock the UpdateWorkervm API call
	client.On("PutJSON", mock.Anything, mock.MatchedBy(func(path string) bool {
		return path == "configuration/v1/worker_vm/123/"
	}), mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updated = true // Mark as updated for subsequent reads
		workerVM := args.Get(3).(*config.WorkerVM)
		*workerVM = *clearAllFieldsMockState
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
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "name", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "hostname", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "domain", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "system_location", "/api/admin/configuration/v1/system_location/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ipv6_gateway", "2001:db8::fe"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "description", "initial description"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "transcoding", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "maintenance_mode_reason", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_cpu_count", "4"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "vm_system_memory", "4096"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_address", "172.16.0.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "secondary_netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "media_priority_weight", "10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "ssh_authorized_keys.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_worker_vm.worker-vm-test", "ssh_authorized_keys.*", "/api/admin/configuration/v1/ssh_authorized_key/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "static_nat_address", "203.0.113.2"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "static_routes.#", "1"),
					resource.TestCheckTypeSetElemAttr("pexip_infinity_worker_vm.worker-vm-test", "static_routes.*", "/api/admin/configuration/v1/static_route/1/"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "tls_certificate", "/api/admin/configuration/v1/tls_certificate/2/"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_location", "test-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_authentication_password", "auth-password1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_community", "public1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_mode", "STANDARD"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_privacy_password", "privacy-password1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_system_contact", "snmpcontact1@domain.com"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "snmp_username", "snmp-user1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "enable_ssh", "ON"),
				),
			},
			{
				Config: test.LoadTestFolder(t, "resource_infinity_worker_vm_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_worker_vm.worker-vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "name", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "hostname", "worker-vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "domain", "updated-value"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "address", "192.168.1.10"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttr("pexip_infinity_worker_vm.worker-vm-test", "system_location", "/api/admin/configuration/v1/system_location/200/"),
				),
			},
		},
	})
}
