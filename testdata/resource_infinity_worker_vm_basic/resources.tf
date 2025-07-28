resource "pexip_infinity_worker_vm" "worker_vm-test" {
  name = "worker_vm-test"
  hostname = "worker_vm-test"
  domain = "test-value"
  address = "192.168.1.10"
  netmask = "255.255.255.0"
  gateway = "192.168.1.1"
  ipv6_address = "2001:db8::1"
  ipv6_gateway = "2001:db8::fe"
  node_type = "conferencing"
  transcoding = true
  password = "test-value"
  maintenance_mode = true
  maintenance_mode_reason = "test-value"
  system_location = "test-value"
}