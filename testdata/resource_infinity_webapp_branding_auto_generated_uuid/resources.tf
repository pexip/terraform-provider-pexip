/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_branding" "webapp_branding-autogen" {
  name          = "webapp_branding-autogen"
  description   = "Test Auto-generated UUID"
  webapp_type   = "webapp2"
  is_default    = false
  branding_file = "${path.module}/webapp2-brand.zip"
  # Note: uuid is not specified, should be auto-generated
}
