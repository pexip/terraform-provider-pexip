terraform {
  required_providers {
    pexip = {
      source  = "pexip"
      version = ">= 1.0.0"
    }
  }
}

provider "pexip" {
  address  = "https://infinity.example.com"
  username = "admin"
  password = "admin_password"
}

resource "pexip_infinity_node" "node" {
  name                    = "test-node-1"
  hostname                = "test-node-1"
  address                 = "192.168.1.100"
  netmask                 = "255.255.255.0"
  domain                  = "example.com"
  gateway                 = "192.168.1.1"
  password                = "test_password"
  node_type               = "worker"
  system_location         = "Test Location"
  maintenance_mode        = false
  maintenance_mode_reason = ""
  transcoding             = true
  vm_cpu_count            = 4
  vm_system_memory        = 8192
}