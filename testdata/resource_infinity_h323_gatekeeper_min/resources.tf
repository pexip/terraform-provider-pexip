/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_h323_gatekeeper" "tf-test-h323-gatekeeper" {
  name    = "tf-test-h323-gatekeeper"
  address = "192.168.1.101"
}
