/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_sip_credential" "sip_credential-test" {
  realm    = "test-value"
  username = "sip_credential-test"
  password = "test-value"
}