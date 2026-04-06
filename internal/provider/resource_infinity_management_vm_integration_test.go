//go:build integration

/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package provider

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/pexip/go-infinity-sdk/v38"

	"github.com/pexip/terraform-provider-pexip/internal/test"
)

func TestInfinityManagementVMIntegration(t *testing.T) {
	_ = os.Setenv("TF_ACC", "1")

	client, err := infinity.New(
		infinity.WithBaseURL(test.INFINITY_BASE_URL),
		infinity.WithBasicAuth(test.INFINITY_USERNAME, test.INFINITY_PASSWORD),
		infinity.WithMaxRetries(2),
		infinity.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // We need this because default certificate is not trusted
				MinVersion:         tls.VersionTLS12,
			},
			MaxIdleConns:        30,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     60 * time.Second,
		}),
	)
	require.NoError(t, err)

	testInfinityManagementVMIntegration(t, client)
}

func testInfinityManagementVMIntegration(t *testing.T, client InfinityClient) {
	// pexip_infinity_management_vm is a singleton — its Delete operation sends a PATCH to reset
	// all fields to API defaults, clearing any references to related resources. Once those
	// references are cleared, the related resources can be deleted normally in the destroy step.
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: getTestProtoV5ProviderFactories(client),
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source: "hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Apply full configuration with all related resources.
			{
				Config: test.LoadTestFolder(t, "resource_infinity_management_vm_full_integration"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "id"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "name", "tf-test-management-vm"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "description", "Integration test management VM"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "mtu", "1400"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "enable_ssh", "ON"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ssh_authorized_keys_use_cloud", "false"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_mode", "AUTHPRIV"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_contact", "admin@example.com"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "snmp_system_location", "tf-test-datacenter"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "dns_servers.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ntp_servers.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "syslog_servers.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "ssh_authorized_keys.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "event_sinks.#", "1"),
					resource.TestCheckResourceAttr("pexip_infinity_management_vm.management_vm-test", "static_routes.#", "1"),
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "http_proxy"),
					resource.TestCheckResourceAttrSet("pexip_infinity_management_vm.management_vm-test", "snmp_network_management_system"),
				),
			},
			// Step 2: Destroy — triggers Delete which resets all fields to API defaults and
			// clears references to related resources, allowing them to be deleted cleanly.
			{
				Config:  test.LoadTestFolder(t, "resource_infinity_management_vm_full_integration"),
				Destroy: true,
			},
		},
	})
}
