resource "pexip_infinity_snmp_network_management_system" "snmp_network_management_system-test" {
  name                = "snmp_network_management_system-test"
  description         = "Test SnmpNetworkManagementSystem"
  address             = "192.168.1.100"
  port                = 162
  snmp_trap_community = "test-value"
}