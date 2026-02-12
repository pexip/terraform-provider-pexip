/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_system_location" "test" {
  name = "tf-test provider system location"
}

# Keep SSH key and static route resources to avoid deletion order issues
# They exist but are not referenced by worker_vm in minimal config
resource "pexip_infinity_ssh_authorized_key" "test" {
  keytype = "ssh-rsa"
  key     = "AAAAB3NzaC1yc2EAAAADAQABAAABgQC7"
  comment = "tf-test SSH key for worker VM"
}

resource "pexip_infinity_static_route" "test" {
  name    = "tf-test static route"
  address = "10.0.0.0"
  prefix  = 24
  gateway = "192.168.1.254"
}

resource "pexip_infinity_worker_vm" "worker-vm-test" {
  # Required fields only - all optional fields cleared
  name            = "tf-test-min_worker-vm"
  hostname        = "worker-vm-test"
  domain          = "test-value" // Keep same - has RequiresReplace
  address         = "192.168.1.10"
  netmask         = "255.255.255.0"
  gateway         = "192.168.1.1"
  system_location = pexip_infinity_system_location.test.id
  password        = "password-initial"

  // Keep RequiresReplace fields to avoid resource replacement
  ipv6_address = "2001:db8::1"
  ipv6_gateway = "2001:db8::fe"

  // All other optional fields are removed to test clearing behavior
}