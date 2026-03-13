/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_gateway_routing_rule" "tf-test-gateway-routing-rule" {
  name         = "tf-test-gateway-routing-rule"
  match_string = ".*@example.com"
  priority     = 100
}
