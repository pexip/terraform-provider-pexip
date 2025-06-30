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

resource "pexip_infinity_dns_server" "cloudflare-dns" {
  address     = "1.1.1.1"
  description = "Cloudflare DNS - updated"
}
