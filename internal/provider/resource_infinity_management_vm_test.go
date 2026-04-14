/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/pexip/go-infinity-sdk/v38/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityManagementVM(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	mockState := &config.ManagementVM{
		ID:                        1,
		ResourceURI:               "/api/admin/configuration/v1/management_vm/1/",
		Name:                      "management_vm-test",
		Description:               "",
		Address:                   "192.168.1.100",
		Netmask:                   "255.255.255.0",
		Gateway:                   "192.168.1.1",
		Hostname:                  "management_vm-test",
		Domain:                    "example.com",
		MTU:                       1500,
		EnableSSH:                 "GLOBAL",
		SSHAuthorizedKeysUseCloud: true,
		SNMPMode:                  "DISABLED",
		SNMPCommunity:             "public",
		SNMPSystemContact:         "admin@domain.com",
		SNMPSystemLocation:        "Virtual machine",
		Primary:                   true,
	}

	// Delete mock — registered first so it takes priority over the general mock.
	// Fingerprinted by Name == "" (delete does not set the name field).
	client.On("PatchJSON", mock.Anything, "configuration/v1/management_vm/1/",
		mock.MatchedBy(func(req *config.ManagementVMUpdateRequest) bool {
			return req.Name == ""
		}), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.ManagementVMUpdateRequest)
		result := args.Get(3).(*config.ManagementVM)

		assert.Equal(t, "", req.Description)
		assert.Equal(t, 1500, req.MTU)
		assert.Equal(t, "GLOBAL", req.EnableSSH)
		assert.True(t, req.SSHAuthorizedKeysUseCloud)
		assert.Equal(t, "DISABLED", req.SNMPMode)
		assert.Equal(t, "public", req.SNMPCommunity)
		assert.Equal(t, "admin@domain.com", req.SNMPSystemContact)
		assert.Equal(t, "Virtual machine", req.SNMPSystemLocation)
		assert.Nil(t, req.TLSCertificate)
		assert.Nil(t, req.SNMPNetworkManagementSystem)
		assert.False(t, req.Initializing)
		assert.Equal(t, []string{}, req.DNSServers)
		assert.Equal(t, []string{}, req.NTPServers)
		assert.Equal(t, []string{}, req.SyslogServers)
		assert.Equal(t, []string{}, req.SSHAuthorizedKeys)
		assert.Equal(t, []string{}, req.StaticRoutes)
		assert.Equal(t, []string{}, req.EventSinks)

		// Reset mockState to defaults
		mockState.Description = ""
		mockState.MTU = 1500
		mockState.StaticNATAddress = nil
		mockState.HTTPProxy = nil
		mockState.TLSCertificate = nil
		mockState.EnableSSH = "GLOBAL"
		mockState.SSHAuthorizedKeysUseCloud = true
		mockState.SNMPMode = "DISABLED"
		mockState.SNMPCommunity = "public"
		mockState.SNMPUsername = ""
		mockState.SNMPAuthenticationPassword = ""
		mockState.SNMPPrivacyPassword = ""
		mockState.SNMPSystemContact = "admin@domain.com"
		mockState.SNMPSystemLocation = "Virtual machine"
		mockState.SNMPNetworkManagementSystem = nil
		mockState.Initializing = false
		mockState.StaticRoutes = nil
		mockState.EventSinks = nil
		mockState.SSHAuthorizedKeys = nil

		*result = *mockState
	}).Once()

	// General PatchJSON mock — handles all create and update calls.
	client.On("PatchJSON", mock.Anything, "configuration/v1/management_vm/1/",
		mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		req := args.Get(2).(*config.ManagementVMUpdateRequest)
		result := args.Get(3).(*config.ManagementVM)

		if req.Name != "" {
			mockState.Name = req.Name
		}
		mockState.Description = req.Description
		mockState.IPV6Address = req.IPV6Address
		mockState.IPV6Gateway = req.IPV6Gateway
		if req.MTU != 0 {
			mockState.MTU = req.MTU
		}
		mockState.StaticNATAddress = req.StaticNATAddress
		mockState.HTTPProxy = req.HTTPProxy
		mockState.TLSCertificate = req.TLSCertificate
		if req.EnableSSH != "" {
			mockState.EnableSSH = req.EnableSSH
		}
		mockState.SSHAuthorizedKeysUseCloud = req.SSHAuthorizedKeysUseCloud
		if req.SNMPMode != "" {
			mockState.SNMPMode = req.SNMPMode
		}
		if req.SNMPCommunity != "" {
			mockState.SNMPCommunity = req.SNMPCommunity
		}
		mockState.SNMPUsername = req.SNMPUsername
		mockState.SNMPAuthenticationPassword = req.SNMPAuthenticationPassword
		mockState.SNMPPrivacyPassword = req.SNMPPrivacyPassword
		if req.SNMPSystemContact != "" {
			mockState.SNMPSystemContact = req.SNMPSystemContact
		}
		if req.SNMPSystemLocation != "" {
			mockState.SNMPSystemLocation = req.SNMPSystemLocation
		}
		mockState.SNMPNetworkManagementSystem = req.SNMPNetworkManagementSystem
		mockState.Initializing = req.Initializing

		*result = *mockState
	}).Maybe()

	// GetJSON mock — returns current mockState for all Read operations.
	client.On(
		"GetJSON",
		mock.Anything,
		"configuration/v1/management_vm/1/",
		mock.Anything,
		mock.AnythingOfType("*config.ManagementVM"),
	).Return(nil).Run(func(args mock.Arguments) {
		mv := args.Get(3).(*config.ManagementVM)
		*mv = *mockState
	}).Maybe()

	testInfinityManagementVM(t, client)
}

func testInfinityManagementVM(t *testing.T, client InfinityClient) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			// Step 1: Set all configurable values
			{
				Config: test.LoadTestFolder(t, "resource_infinity_management_vm_full"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "id"),
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "resource_id"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "name", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "description", "Test ManagementVM"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "address", "192.168.1.100"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "gateway", "192.168.1.1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "hostname", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "domain", "example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "mtu", "1400"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ipv6_address", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ipv6_gateway", "2001:db8::1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "static_nat_address", "192.0.2.1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "http_proxy", "http://proxy.example.com:8080"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "tls_certificate", "test-certificate"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "enable_ssh", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ssh_authorized_keys_use_cloud", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_mode", "AUTHPRIV"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_community", "public"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_username", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_authentication_password", "test-auth-pass"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_privacy_password", "test-priv-pass"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_contact", "admin@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_location", "datacenter"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_network_management_system", "192.168.1.200"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "primary", "true"),
				),
			},
			// Step 2: Destroy — triggers the delete which resets all fields to defaults.
			{
				Destroy: true,
				Config:  test.LoadTestFolder(t, "resource_infinity_management_vm_full"),
			},
			// Step 3: Re-apply min config and verify all fields are at their defaults,
			// confirming the delete actually reset the API state.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_management_vm_min"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "name", "management_vm-test"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "description", ""),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "mtu", "1500"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "enable_ssh", "GLOBAL"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ssh_authorized_keys_use_cloud", "true"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_mode", "DISABLED"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_community", "public"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_username", ""),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_contact", "admin@domain.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_location", "Virtual machine"),
					resource.TestCheckNoResourceAttr("pexip_infinity_management_vm.management_vm-test", "ipv6_address"),
					resource.TestCheckNoResourceAttr("pexip_infinity_management_vm.management_vm-test", "ipv6_gateway"),
					resource.TestCheckNoResourceAttr("pexip_infinity_management_vm.management_vm-test", "static_nat_address"),
					resource.TestCheckNoResourceAttr("pexip_infinity_management_vm.management_vm-test", "http_proxy"),
					resource.TestCheckNoResourceAttr("pexip_infinity_management_vm.management_vm-test", "tls_certificate"),
					resource.TestCheckNoResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_network_management_system"),
				),
			},
		},
	})
}

func TestInfinityManagementVMValidation(t *testing.T) {
	t.Parallel()
	_ = os.Setenv("TF_ACC", "1")

	client := infinity.NewClientMock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: getTestProtoV6ProviderFactories(client),
		Steps: []resource.TestStep{
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name = ""
}
`,
				ExpectError: regexp.MustCompile(`string length must be at least 1`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name = "` + strings.Repeat("a", 33) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 32`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name        = "management_vm-test"
  description = "` + strings.Repeat("a", 251) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 250`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name             = "management_vm-test"
  alternative_fqdn = "` + strings.Repeat("a", 256) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 255`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name = "management_vm-test"
  mtu  = 511
}
`,
				ExpectError: regexp.MustCompile(`value must be between 512 and 1500`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name = "management_vm-test"
  mtu  = 1501
}
`,
				ExpectError: regexp.MustCompile(`value must be between 512 and 1500`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name       = "management_vm-test"
  enable_ssh = "INVALID"
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name      = "management_vm-test"
  snmp_mode = "INVALID"
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name           = "management_vm-test"
  snmp_community = "` + strings.Repeat("a", 17) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 16`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name          = "management_vm-test"
  snmp_username = "` + strings.Repeat("a", 101) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 100`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name                        = "management_vm-test"
  snmp_authentication_password = "1234567"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at least 8`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name                        = "management_vm-test"
  snmp_authentication_password = "` + strings.Repeat("a", 101) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 100`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name                 = "management_vm-test"
  snmp_privacy_password = "1234567"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at least 8`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name                 = "management_vm-test"
  snmp_privacy_password = "` + strings.Repeat("a", 101) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 100`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name                = "management_vm-test"
  snmp_system_contact = "` + strings.Repeat("a", 71) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 70`),
			},
			{
				Config: `
resource "pexip_infinity_management_vm" "management_vm-test" {
  name                 = "management_vm-test"
  snmp_system_location = "` + strings.Repeat("a", 71) + `"
}
`,
				ExpectError: regexp.MustCompile(`string length must be at most 70`),
			},
		},
	})
}
