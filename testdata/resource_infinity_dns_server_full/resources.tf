/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_dns_server" "tf-test-dns" {
  address     = "4.2.2.2"
  description = "tf-test Level 3 DNS Server"
}
