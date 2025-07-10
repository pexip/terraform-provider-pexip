resource "pexip_infinity_node" "node" {
  name                    = "test-node-1"
  hostname                = "test-node-1"
  address                 = "192.168.1.100"
  netmask                 = "255.255.255.0"
  domain                  = "pexip.com" # Updated domain from example.com to pexip.com
  gateway                 = "192.168.1.1"
  password                = "password123"
  node_type               = "CONFERENCING"
  system_location         = "Test Location"
  maintenance_mode        = false
  maintenance_mode_reason = ""
  transcoding             = true
  vm_cpu_count            = 4
  vm_system_memory        = 8192
}