/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ntp_server" "tf-test-ntp" {
  address     = "1.europe.pool.ntp.org"
  description = "tf-test NTP Server Description"
}
