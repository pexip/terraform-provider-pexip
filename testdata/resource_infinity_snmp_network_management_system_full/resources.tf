/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_snmp_network_management_system" "tf-test-snmp-nms" {
  name                = "tf-test-snmp-nms"
  description         = "tf-test SNMP NMS Description"
  address             = "192.168.1.100"
  port                = 162
  snmp_trap_community = "tf-test-comm"
}
