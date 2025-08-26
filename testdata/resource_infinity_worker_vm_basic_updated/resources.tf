/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_worker_vm" "worker-vm-test" {
  name                    = "worker-vm-test"
  hostname                = "worker-vm-test"
  domain                  = "updated-value"
  address                 = "192.168.1.10"
  netmask                 = "255.255.255.0"
  gateway                 = "192.168.1.1"
  system_location         = "/api/admin/configuration/v1/system_location/2/"
  ipv6_address            = "2001:db8::2"
  ipv6_gateway            = "2001:db8::ff"
  description             = "updated description"
  snmp_authentication_password = "auth-password2"
  snmp_community          = "public"
  snmp_mode               = "AUTHPRIV"
  snmp_privacy_password       = "privacy-password2"
  snmp_system_contact      = "snmpcontact2@domain.com"
  snmp_system_location     = "eu2"
  snmp_username            = "snmp-user2"
  enable_ssh             = "OFF"
  enable_distributed_database = false
  media_priority_weight  = 100
  cloud_bursting         = true
  ssh_authorized_keys_use_cloud = false
  tls_certificate       = "/api/admin/configuration/v1/tls_certificate/1/"
  static_nat_address = "203.0.113.2"
}