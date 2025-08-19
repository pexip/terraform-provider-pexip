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
      version = "3.3.2"
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
  auth_url    = var.auth_url
  user_name   = var.openstack_username
  password    = var.openstack_password
  tenant_name = var.openstack_tenant_name
  domain_name = var.openstack_domain_name
}

provider "pexip" {
  //address  = "https://${local.manager_hostname}.${local.domain}"
  address  = "https://${module.openstack-infinity-manager.mgr-public-ip}"
  username = var.infinity_username
  password = var.infinity_password
  insecure = true
}
