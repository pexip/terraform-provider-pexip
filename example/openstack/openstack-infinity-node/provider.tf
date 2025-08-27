/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

terraform {
  required_version = ">= 1.0"
  required_providers {
    pexip = {
      source  = "pexip/pexip"
      version = "0.9.0"
    }
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "3.3.2"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 3.0"
    }
  }
}