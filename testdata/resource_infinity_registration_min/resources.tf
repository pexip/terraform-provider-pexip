/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_registration" "registration-test" {
  refresh_strategy              = "adaptive"
  route_via_registrar           = true
  enable_push_notifications     = false
  enable_google_cloud_messaging = true
}
