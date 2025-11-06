/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ivr_theme" "ivr_theme-test" {
  name    = "ivr_theme-test"
  package = "ivr_theme_test.zip"
}