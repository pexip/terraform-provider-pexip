/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_alias" "scheduled_alias-test" {
  alias              = "updated-value" // Updated value
  alias_number       = 9876543210      // Updated value
  numeric_alias      = "updated-value" // Updated value
  uuid               = "updated-value" // Updated value
  exchange_connector = "updated-value" // Updated value
  is_used            = false           // Updated to false
  ews_item_uid       = "updated-value" // Updated value
}