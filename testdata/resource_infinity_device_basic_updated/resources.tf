/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_device" "device-test" {
  alias                           = "device-test"
  description                     = "Updated Test Device" // Updated description
  username                        = "updateduser"         // Updated username
  password                        = "updatedpass"         // Updated password
  primary_owner_email_address     = "updated@example.com" // Updated email
  enable_sip                      = false                 // Updated to false
  enable_h323                     = true                  // Updated to true
  enable_infinity_connect_non_sso = false                 // Updated to false
  enable_infinity_connect_sso     = false                 // Keep disabled (requires identity provider)
  enable_standard_sso             = false                 // Keep disabled (requires identity provider)
  tag                             = "updated-tag"         // Updated tag
  sync_tag                        = "updated-sync-tag"    // Updated sync tag
}