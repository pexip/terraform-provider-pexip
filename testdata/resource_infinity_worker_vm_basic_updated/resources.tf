/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_worker_vm" "worker-vm-test" {
  name            = "worker-vm-test"
  hostname        = "worker-vm-test"
  domain          = "test-value" // Keep same - has RequiresReplace
  address         = "192.168.1.10"
  netmask         = "255.255.255.0"
  gateway         = "192.168.1.1"
  system_location = "/api/admin/configuration/v1/system_location/1/"
  password        = "password-updated"
  description     = "updated description"

  // Keep RequiresReplace fields to avoid resource replacement
  ipv6_address = "2001:db8::1"
  ipv6_gateway = "2001:db8::fe"

  // Add some optional fields to test update
  maintenance_mode        = false
}