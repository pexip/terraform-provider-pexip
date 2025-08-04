resource "pexip_infinity_management_vm" "management_vm-test" {
  name                           = "management_vm-test"
  description                    = "Updated Test ManagementVM" // Updated description
  address                        = "192.168.1.101"             // Updated address
  netmask                        = "255.255.252.0"             // Updated value
  gateway                        = "192.168.1.254"             // Updated value
  mtu                            = 1400                        // Updated MTU
  hostname                       = "management_vm-test"
  domain                         = "updated.com"                           // Updated value
  alternative_fqdn               = "alt-updated.example.com"               // Updated value
  ipv6_address                   = "2001:db8::2"                           // Updated address
  ipv6_gateway                   = "2001:db8::254"                         // Updated value
  static_nat_address             = "203.0.113.2"                           // Updated address
  http_proxy                     = "http://updated-proxy.example.com:3128" // Updated value
  tls_certificate                = "updated-certificate"                   // Updated value
  enable_ssh                     = "keys_only"                             // Updated value
  ssh_authorized_keys_use_cloud  = false                                   // Updated to false
  secondary_config_passphrase    = "updated-passphrase"                    // Updated value
  snmp_mode                      = "v3"                                    // Updated value
  snmp_community                 = "private"                               // Updated value
  snmp_username                  = "management_vm-test"
  snmp_authentication_password   = "updated-auth-pass"         // Updated value
  snmp_privacy_password          = "updated-priv-pass"         // Updated value
  snmp_system_contact            = "updated-admin@example.com" // Updated value
  snmp_system_location           = "updated-datacenter"        // Updated value
  snmp_network_management_system = "192.168.1.201"             // Updated value
  initializing                   = false                       // Updated to false
}