/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_oauth2_client" "oauth2_client-test" {
  client_name = "oauth2_client-test"
  role        = "test-value"
}