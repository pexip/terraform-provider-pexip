/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

terraform {
  required_providers {
    pexip = {
      source  = "pexip.com/pexip/pexip"
      version = "0.0.1"
    }
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.54.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.7.1"
    }
    local = {
      source  = "hashicorp/local"
      version = "2.4.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.2.1"
    }
  }
}

provider "tls" {}

provider "openstack" {
  auth_url = "https://bf-ososl-horizon.osl.basefarm.net:5000/v3"
  user_name = "massel"
  password = "iega5ua0o-iSh4sh" # export OS_PASSWORD=your-password
  tenant_name = "pexi_dev" #project
  domain_name = "pexi"
}

provider "pexip" {
  //address  = "https://${local.manager_hostname}.${local.domain}"
  address  = "https://${module.openstack-infinity-manager.mgr-public-ip}"
  username = var.infinity_username
  password = var.infinity_password
  insecure = true
}
