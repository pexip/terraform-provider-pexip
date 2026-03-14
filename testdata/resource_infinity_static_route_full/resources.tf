/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */
resource "pexip_infinity_static_route" "tf-test-static-route" {
  name    = "tf-test-static-route-updated"
  address = "10.0.0.0"
  prefix  = 16
  gateway = "10.0.0.1"
}
