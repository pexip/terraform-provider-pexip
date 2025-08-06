/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_recurring_conference" "recurring_conference-test" {
  conference      = "updated-value" // Updated value
  current_index   = 2               // Updated value
  ews_item_id     = "updated-value" // Updated value
  is_depleted     = false           // Updated to false
  subject         = "updated-value" // Updated value
  scheduled_alias = "updated-value" // Updated value
}