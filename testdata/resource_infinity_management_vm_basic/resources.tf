/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_management_vm" "management_vm-test" {
  name                           = "management_vm-test"
  description                    = "Test ManagementVM"
  address                        = "192.168.1.100"
  netmask                        = "255.255.255.0"
  gateway                        = "192.168.1.1"
  mtu                            = 1500
  hostname                       = "management_vm-test"
  domain                         = "example.com"
  alternative_fqdn               = "alt.example.com"
  ipv6_address                   = "2001:db8::1"
  ipv6_gateway                   = "2001:db8::1"
  static_nat_address             = "203.0.113.1"
  http_proxy                     = "http://proxy.example.com:8080"
  tls_certificate                = "test-certificate"
  enable_ssh                     = "yes"
  ssh_authorized_keys_use_cloud  = true
  secondary_config_passphrase    = "test-passphrase"
  snmp_mode                      = "v1v2c"
  snmp_community                 = "public"
  snmp_username                  = "management_vm-test"
  snmp_authentication_password   = "test-auth-pass"
  snmp_privacy_password          = "test-priv-pass"
  snmp_system_contact            = "admin@example.com"
  snmp_system_location           = "datacenter"
  snmp_network_management_system = "192.168.1.200"
  initializing                   = true
}