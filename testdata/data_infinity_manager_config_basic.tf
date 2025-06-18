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

data "pexip_infinity_manager_config" "master" {
  hostname              = "test-mgr1"
  domain                = "dev.vcops.tech"
  ip                    = "10.0.0.40"
  mask                  = "255.255.255.0"
  gw                    = "10.0.0.1"
  dns                   = "1.1.1.1"
  ntp                   = "pool.ntp.org"
  user                  = "admin"
  pass                  = "admin_password"
  admin_password        = "admin_password"
  error_reports         = false
  enable_analytics      = false
  contact_email_address = "vcops@pexip.com"
}