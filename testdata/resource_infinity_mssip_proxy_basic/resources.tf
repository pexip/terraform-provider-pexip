/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mssip_proxy" "mssip-proxy-test" {
  name        = "mssip-proxy-test"
  description = "Test MSSIP proxy"
  address     = "test-mssip-proxy.dev.pexip.network"
  port        = 5060
  transport   = "tcp"
}