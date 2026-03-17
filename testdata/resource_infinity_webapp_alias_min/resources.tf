/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_alias" "tf-test-webapp-alias" {
  slug        = "tf-test-alias"
  webapp_type = "webapp1"
}
