/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ntp_server" "ntp-2" {
  address     = "2.pool.ntp.org"
  description = "NTP server 2"
}
