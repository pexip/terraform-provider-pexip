/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_licence" "licence-test" {
  entitlement_id = "test-value"
  offline_mode   = true
}