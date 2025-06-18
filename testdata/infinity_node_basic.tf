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
  name   = "test-node-1"
  config = "<config from manager>"
}