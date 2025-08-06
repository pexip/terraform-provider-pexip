/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_static_route" "static_route-test" {
  name    = "static_route-test"
  address = "192.168.1.0"
  prefix  = 24
  gateway = "192.168.1.1"
}