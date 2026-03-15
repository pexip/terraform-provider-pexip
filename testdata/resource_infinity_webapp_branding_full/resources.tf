/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_branding" "webapp_branding-test" {
  name          = "tf-test-webapp-branding-full"
  description   = "tf-test webapp branding description"
  webapp_type   = "webapp3"
  branding_file = "${path.module}/webapp3-brand.zip"
}
