/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_sip_proxy" "sip-proxy-test" {
  name        = "tf-test-sip-proxy-full"
  description = "Full configuration test SIP proxy"
  address     = "sip.pexvclab.com"
  port        = 5061
  transport   = "tls"
}
