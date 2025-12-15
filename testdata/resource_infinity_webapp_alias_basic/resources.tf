/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_webapp_alias" "webapp_alias-test" {
  slug        = "test-alias"
  description = "Test WebappAlias"
  webapp_type = "webapp1"
  is_enabled  = true
}