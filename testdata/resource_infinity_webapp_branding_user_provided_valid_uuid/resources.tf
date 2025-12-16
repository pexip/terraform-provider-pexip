/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_branding" "webapp_branding-valid-uuid" {
  name          = "webapp_branding-valid-uuid"
  description   = "Test Valid UUID"
  uuid          = "550e8400-e29b-41d4-a716-446655440000"
  webapp_type   = "webapp3"
  is_default    = true
  branding_file = "${path.module}/webapp2-brand.zip"
}
