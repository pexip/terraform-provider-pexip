/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_worker_vm" "worker-vm-test" {
  name            = "worker-vm-test"
  hostname        = "worker-vm-test"
  domain          = "updated-value"
  address         = "192.168.1.10"
  netmask         = "255.255.255.0"
  gateway         = "192.168.1.1"
  system_location = "/api/admin/configuration/v1/system_location/200/"
}