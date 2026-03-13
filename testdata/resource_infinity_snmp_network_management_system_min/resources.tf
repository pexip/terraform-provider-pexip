/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_snmp_network_management_system" "tf-test-snmp-nms" {
  name    = "tf-test-snmp-nms"
  address = "192.168.1.100"
}
