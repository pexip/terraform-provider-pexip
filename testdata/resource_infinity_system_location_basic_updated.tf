terraform {
  required_providers {
    pexip = {
      source  = "pexip"
      version = ">= 1.0.0"
    }
  }
}

provider "pexip" {
  address  = "https://dev-manager.dev.pexip.network"
  username = "admin"
  password = "admin"
  insecure = true
}

resource "pexip_infinity_system_location" "main-location" {
  name        = "main"
  description = "Main location for Pexip Infinity System - updated" # Updated description
  mtu         = 1460
  dns_servers = ["/api/admin/configuration/v1/dns_server/1/"] # Updated DNS servers
  ntp_servers = ["/api/admin/configuration/v1/ntp_server/1/"]
}
