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

resource "pexip_infinity_ntp_server" "ntp-2" {
  address     = "2.pool.ntp.org"
  description = "NTP server 2"
}
