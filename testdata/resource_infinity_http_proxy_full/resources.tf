/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_http_proxy" "tf-test-http-proxy" {
  name     = "tf-test-http-proxy"
  address  = "proxy.example.com"
  port     = 8081
  protocol = "http"
  username = "tf-test-user"
  password = "tf-test-password"
}
