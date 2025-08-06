/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_sip_credential" "sip_credential-test" {
  realm    = "updated-value" // Updated value
  username = "sip_credential-test"
  password = "updated-value" // Updated value
}