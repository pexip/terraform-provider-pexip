/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ivr_theme" "tf-test-theme" {
  name = "tf-test-theme"
}

resource "pexip_infinity_global_configuration" "global_configuration-test" {
  default_theme             = pexip_infinity_ivr_theme.tf-test-theme.id
  logon_banner              = "tf-test-logon-banner"
  site_banner               = "tf-test-site-banner"
  site_banner_bg            = "#123456"
  site_banner_fg            = "#ffffff"
  contact_email_address     = "tf-test@example.com"
  enable_analytics          = true
  enable_breakout_rooms     = true
  enable_clock              = true
  enable_legacy_dialout_api = true
  guests_only_timeout       = 120
}
