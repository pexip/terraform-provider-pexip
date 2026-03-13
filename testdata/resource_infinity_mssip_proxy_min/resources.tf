/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mssip_proxy" "tf-test-mssip-proxy" {
  name    = "tf-test-mssip-proxy"
  address = "mssip-proxy-min.example.com"
}
