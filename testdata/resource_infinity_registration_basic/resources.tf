/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_registration" "registration-test" {
  enable                        = true
  refresh_strategy              = "adaptive"
  adaptive_min_refresh          = 300
  adaptive_max_refresh          = 600
  natted_min_refresh            = 60
  natted_max_refresh            = 120
  route_via_registrar           = true
  enable_push_notifications     = true
  enable_google_cloud_messaging = true
  push_token                    = "test-value"
}