resource "pexip_infinity_worker_vm" "worker_vm-test" {
  name                    = "worker_vm-test"
  hostname                = "worker_vm-test"
  domain                  = "updated-value" // Updated value
  address                 = "192.168.1.20"  // Updated address
  netmask                 = "255.255.255.0" // Updated value
  gateway                 = "192.168.1.1"   // Updated value
  ipv6_address            = "2001:db8::2"   // Updated address
  ipv6_gateway            = "2001:db8::fe"  // Updated value
  node_type               = "proxying"      // Updated value
  transcoding             = false           // Updated to false
  password                = "updated-value" // Updated value
  maintenance_mode        = false           // Updated to false
  maintenance_mode_reason = "updated-value" // Updated value
  system_location         = "updated-value" // Updated value
}