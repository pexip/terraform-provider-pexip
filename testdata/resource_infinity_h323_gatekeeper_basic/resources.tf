/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_h323_gatekeeper" "h323_gatekeeper-test" {
  name        = "h323_gatekeeper-test"
  description = "Test H323Gatekeeper"
  address     = "192.168.1.100"
}