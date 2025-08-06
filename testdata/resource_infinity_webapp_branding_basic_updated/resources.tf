/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_branding" "webapp_branding-test" {
  name          = "webapp_branding-test"
  description   = "Updated Test WebappBranding" // Updated description
  uuid          = "updated-value"               // Updated value
  webapp_type   = "management"                  // Updated value
  is_default    = false                         // Updated to false
  branding_file = "updated-value"               // Updated value
}