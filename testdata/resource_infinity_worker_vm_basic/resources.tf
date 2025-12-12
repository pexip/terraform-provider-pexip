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
  tls_certificate               = "/api/admin/configuration/v1/tls_certificate/2/"
  description                   = "initial description"
  ipv6_address                  = "2001:db8::1"
  ipv6_gateway                  = "2001:db8::fe"
  node_type                     = "CONFERENCING"
  deployment_type               = "MANUAL-PROVISION-ONLY"
  transcoding                   = true
  password                      = "password-initial"
  maintenance_mode              = true
  maintenance_mode_reason       = "test-value"
  vm_cpu_count                  = 4
  vm_system_memory              = 4096
  secondary_address             = "172.16.0.10"
  secondary_netmask             = "255.255.255.0"
  media_priority_weight         = 10
  ssh_authorized_keys_use_cloud = true
  static_nat_address            = "203.0.113.2"
  snmp_authentication_password  = "auth-password1"
  snmp_community                = "public1"
  snmp_mode                     = "STANDARD"
  snmp_privacy_password         = "privacy-password1"
  snmp_system_contact           = "snmpcontact1@domain.com"
  snmp_system_location          = "test-value"
  snmp_username                 = "snmp-user1"
  enable_ssh                    = "ON"
  enable_distributed_database   = false
  service_manager               = true
  service_policy                = true
  signalling                    = true
  managed                       = false
}