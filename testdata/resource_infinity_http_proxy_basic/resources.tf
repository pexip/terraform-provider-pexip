/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_http_proxy" "http_proxy-test" {
  name     = "http_proxy-test"
  address  = "test-server.example.com"
  port     = 8080
  protocol = "https"
  username = "http_proxy-test"
  password = "test-value"
}