/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

terraform {
  required_providers {
    pexip = {
      source  = "pexip"
      version = "0.0.1"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

provider "pexip" {
  address  = "https://dev-manager.dev.pexip.network"
  username = "admin"
  password = "admin"
  insecure = true
}
