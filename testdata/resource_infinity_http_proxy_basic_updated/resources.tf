/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_http_proxy" "http_proxy-test" {
  name     = "http_proxy-test"
  address  = "updated-server.example.com" // Updated address
  port     = 8081                         // Updated port
  protocol = "http"                       // Updated value
  username = "http_proxy-test"
  password = "updated-value" // Updated value
}