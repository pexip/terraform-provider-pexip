/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_alias" "scheduled_alias-test" {
  alias              = "test-value"
  alias_number       = 1234567890
  numeric_alias      = "test-value"
  uuid               = "test-value"
  exchange_connector = "test-value"
  is_used            = true
  ews_item_uid       = "test-value"
}