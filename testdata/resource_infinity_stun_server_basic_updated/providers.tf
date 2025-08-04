terraform {
  required_providers {
    pexip = {
      source  = "pexip"
      version = "0.0.1"
    }
  }
}

provider "pexip" {
  address  = "https://dev-manager.dev.pexip.network"
  username = "admin"
  password = "admin"
  insecure = true
}
