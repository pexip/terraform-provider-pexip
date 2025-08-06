/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_registration" "registration-test" {
  enable                        = false           // Updated to false
  refresh_strategy              = "maximum"       // Updated value
  route_via_registrar           = false           // Updated to false
  enable_push_notifications     = false           // Updated to false
  enable_google_cloud_messaging = false           // Updated to false
  push_token                    = "updated-value" // Updated value
}