/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_dns_server" "cloudflare-dns" {
  address     = "192.168.1.50"
  description = "Cloudflare DNS - updated"
}
