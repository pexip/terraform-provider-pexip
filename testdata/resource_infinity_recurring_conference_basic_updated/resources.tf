/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "test-conference" {
  name         = "test-conference"
  description  = "Test Conference"
  service_type = "conference"
}

resource "pexip_infinity_recurring_conference" "recurring_conference-test" {
  conference    = pexip_infinity_conference.test-conference.id
  current_index = 2                                   // Updated value
  ews_item_id   = "updated-ews-item-id"               // Updated value
  is_depleted   = true                                // Updated to true
  subject       = "Updated Test Recurring Conference" // Updated value

  depends_on = [
    pexip_infinity_conference.test-conference
  ]
}