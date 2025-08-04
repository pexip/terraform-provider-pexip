resource "pexip_infinity_snmp_network_management_system" "snmp_network_management_system-test" {
  name                = "snmp_network_management_system-test"
  description         = "Updated Test SnmpNetworkManagementSystem" // Updated description
  address             = "192.168.1.200"                            // Updated address
  port                = 161                                        // Updated port
  snmp_trap_community = "updated-value"                            // Updated value
}