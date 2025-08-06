/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_alias" "webapp_alias-test" {
  slug        = "updated-value"            // Updated value
  description = "Updated Test WebappAlias" // Updated description
  webapp_type = "management"               // Updated value
  is_enabled  = false                      // Updated to false
  bundle      = "updated-value"            // Updated value
  branding    = "updated-value"            // Updated value
}