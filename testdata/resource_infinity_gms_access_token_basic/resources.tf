/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_gms_access_token" "gms_access_token-test" {
  name  = "gms_access_token-test"
  token = "test-value"
}