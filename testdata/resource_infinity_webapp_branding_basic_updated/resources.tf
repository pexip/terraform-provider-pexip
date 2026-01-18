/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_branding" "webapp_branding-test" {
  name          = "webapp_branding-test"
  description   = "Updated Test WebappBranding"          // Updated description
  webapp_type   = "webapp2"
  branding_file = "${path.module}/webapp2-brand.zip"
}