/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_branding" "webapp_branding-test" {
  name          = "webapp_branding-test"
  description   = "Test WebappBranding"
  uuid          = "test-value"
  webapp_type   = "webapp1"
  is_default    = true
  branding_file = "test-value"
}