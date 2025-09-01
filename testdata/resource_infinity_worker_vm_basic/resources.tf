/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_worker_vm" "worker-vm-test" {
  name                          = "worker-vm-test"
  hostname                      = "worker-vm-test"
  domain                        = "test-value"
  address                       = "192.168.1.10"
  netmask                       = "255.255.255.0"
  gateway                       = "192.168.1.1"
  system_location               = "/api/admin/configuration/v1/system_location/1/"
  ipv6_address                  = "2001:db8::1"
  ipv6_gateway                  = "2001:db8::fe"
  description                   = "initial description"
  snmp_system_contact           = "snmpcontact1@domain.com"
  snmp_system_location          = "eu1"
  snmp_username                 = "snmp-user1"
  snmp_authentication_password  = "auth-password1"
  snmp_privacy_password         = "privacy-password1"
  snmp_mode                     = "STANDARD"
  enable_ssh                    = "ON"
  enable_distributed_database   = true
  media_priority_weight         = 0
  cloud_bursting                = false
  ssh_authorized_keys_use_cloud = true
  tls_certificate               = "/api/admin/configuration/v1/tls_certificate/2/"
  static_nat_address            = "203.0.113.2"
  //node_type               = "CONFERENCING"
  //transcoding             = true
  //password                = "test-value"
  maintenance_mode        = true
  //maintenance_mode_reason = "test-value"
}