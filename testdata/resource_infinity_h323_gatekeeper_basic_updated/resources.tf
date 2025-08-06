/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_h323_gatekeeper" "h323_gatekeeper-test" {
  name        = "h323_gatekeeper-test"
  description = "Updated Test H323Gatekeeper" // Updated description
  address     = "192.168.1.200"               // Updated address
}