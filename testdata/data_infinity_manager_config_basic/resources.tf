/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

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