/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_registration" "registration-test" {
  enable                        = false
  refresh_strategy              = "maximum"
  maximum_min_refresh           = 120
  maximum_max_refresh           = 600
  natted_min_refresh            = 120
  natted_max_refresh            = 180
  route_via_registrar           = false
  enable_push_notifications     = true
  enable_google_cloud_messaging = false
  push_token                    = "custom-token"
}
