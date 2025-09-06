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

func TestInfinityManagementVM(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	// Create a mock client and set up expectations
	client := infinity.NewClientMock()

	// Mock the CreateManagementvm API call
	createResponse := &types.PostResponse{
		Body:        []byte(""),
		ResourceURI: "/api/admin/configuration/v1/management_vm/123/",
	}
	client.On("PostWithResponse", mock.Anything, "configuration/v1/management_vm/", mock.Anything, mock.Anything).Return(createResponse, nil)

	// Shared state for mocking
	mockState := &config.ManagementVM{
		ID:                          123,
		ResourceURI:                 "/api/admin/configuration/v1/management_vm/123/",
		Name:                        "management_vm-test",
		Description:                 "Test ManagementVM",
		Address:                     "192.168.1.100",
		Netmask:                     "255.255.255.0",
		Gateway:                     "192.168.1.1",
		MTU:                         1500,
		Hostname:                    "management_vm-test",
		Domain:                      "example.com",
		AlternativeFQDN:             "alt.example.com",
		IPV6Address:                 test.StringPtr("2001:db8::1"),
		IPV6Gateway:                 test.StringPtr("2001:db8::1"),
		HTTPProxy:                   test.StringPtr("http://proxy.example.com:8080"),
		TLSCertificate:              test.StringPtr("test-certificate"),
		EnableSSH:                   "ON",
		SSHAuthorizedKeysUseCloud:   true,
		SecondaryConfigPassphrase:   "test-passphrase",
		SNMPMode:                    "AUTHPRIV",
		SNMPCommunity:               "public",
		SNMPUsername:                "management_vm-test",
		SNMPAuthenticationPassword:  "test-auth-pass",
		SNMPPrivacyPassword:         "test-priv-pass",
		SNMPSystemContact:           "admin@example.com",
		SNMPSystemLocation:          "datacenter",
		SNMPNetworkManagementSystem: test.StringPtr("192.168.1.200"),
		Initializing:                true,
		Primary:                     true,
	}

	// Mock the GetManagementvm API call for Read operations
	client.On(
		"GetJSON",
		mock.Anything,
		"configuration/v1/management_vm/1/",
		mock.Anything,
		mock.AnythingOfType("*config.ManagementVM"), // pointer to config.ManagementVM
	).Return(nil).Run(func(args mock.Arguments) {
		global_configuration := args.Get(3).(*config.ManagementVM)
		*global_configuration = *mockState
	}).Maybe()

	// Mock the UpdateManagementVM API call
	client.On("PatchJSON", mock.Anything, "configuration/v1/management_vm/1/", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updateRequest := args.Get(2).(*config.ManagementVMUpdateRequest)
		management_vm := args.Get(3).(*config.ManagementVM)

		// Update mock state based on request
		management_vm.Name = updateRequest.Name
		management_vm.Description = updateRequest.Description
		management_vm.Address = updateRequest.Address
		management_vm.Netmask = updateRequest.Netmask
		management_vm.Gateway = updateRequest.Gateway
		management_vm.MTU = updateRequest.MTU
		management_vm.Hostname = updateRequest.Hostname
		management_vm.Domain = updateRequest.Domain
		management_vm.AlternativeFQDN = updateRequest.AlternativeFQDN
		management_vm.IPV6Address = updateRequest.IPV6Address
		management_vm.IPV6Gateway = updateRequest.IPV6Gateway
		management_vm.StaticNATAddress = updateRequest.StaticNATAddress
		management_vm.HTTPProxy = updateRequest.HTTPProxy
		management_vm.TLSCertificate = updateRequest.TLSCertificate
		management_vm.EnableSSH = updateRequest.EnableSSH
		management_vm.SSHAuthorizedKeysUseCloud = updateRequest.SSHAuthorizedKeysUseCloud
		management_vm.SecondaryConfigPassphrase = updateRequest.SecondaryConfigPassphrase
		management_vm.SNMPMode = updateRequest.SNMPMode
		management_vm.SNMPCommunity = updateRequest.SNMPCommunity
		management_vm.SNMPUsername = updateRequest.SNMPUsername
		management_vm.SNMPAuthenticationPassword = updateRequest.SNMPAuthenticationPassword
		management_vm.SNMPPrivacyPassword = updateRequest.SNMPPrivacyPassword
		management_vm.SNMPSystemContact = updateRequest.SNMPSystemContact
		management_vm.SNMPSystemLocation = updateRequest.SNMPSystemLocation
		management_vm.SNMPNetworkManagementSystem = updateRequest.SNMPNetworkManagementSystem

		// Return updated state
		*management_vm = *mockState

	}).Maybe()

	testInfinityManagementVM(t, client)
}

func testInfinityManagementVM(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: test.LoadTestFolder(t, "resource_infinity_management_vm_basic_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "name", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "description", "Test ManagementVM"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "mtu", "1500"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "hostname", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "domain", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "alternative_fqdn", "alt.example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ipv6_gateway", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "http_proxy", "http://proxy.example.com:8080"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "tls_certificate", "test-certificate"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "enable_ssh", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ssh_authorized_keys_use_cloud", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "secondary_config_passphrase", "test-passphrase"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_mode", "AUTHPRIV"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_community", "public"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_username", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_authentication_password", "test-auth-pass"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_privacy_password", "test-priv-pass"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_contact", "admin@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_location", "datacenter"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_network_management_system", "192.168.1.200"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "initializing", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "primary", "true"),
				),
			},
			// ManagementVM doesn't support updates, so only test creation/read
		},
	})
}
