/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_management_vm" "management_vm-test" {
  name                           = "management_vm-test"
  description                    = "Test ManagementVM"
  ipv6_address                   = "2001:db8::1"
  ipv6_gateway                   = "2001:db8::1"
  mtu                            = 1500
  http_proxy                     = "http://proxy.example.com:8080"
  tls_certificate                = "test-certificate"
  enable_ssh                     = "ON"
  ssh_authorized_keys_use_cloud  = true
  snmp_mode                      = "AUTHPRIV"
  snmp_community                 = "public"
  snmp_username                  = "management_vm-test"
  snmp_authentication_password   = "test-auth-pass"
  snmp_privacy_password          = "test-priv-pass"
  snmp_system_contact            = "admin@example.com"
  snmp_system_location           = "datacenter"
  snmp_network_management_system = "192.168.1.200"
}
