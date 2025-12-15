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
  current_index = 1
  ews_item_id   = "test-ews-item-id"
  is_depleted   = false
  subject       = "Test Recurring Conference"

  depends_on = [
    pexip_infinity_conference.test-conference
  ]
}