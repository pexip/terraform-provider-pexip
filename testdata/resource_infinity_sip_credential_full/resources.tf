/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_sip_credential" "tf-test-sip-credential" {
  realm    = "tf-test-realm"
  username = "tf-test-sip-credential"
  password = "tf-test-password"
}
