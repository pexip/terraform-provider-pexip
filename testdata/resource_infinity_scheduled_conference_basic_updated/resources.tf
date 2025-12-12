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

resource "pexip_infinity_scheduled_conference" "scheduled_conference-test" {
  conference   = pexip_infinity_conference.test-conference.id
  start_time   = "2024-02-01T14:00:00Z"              // Updated value
  end_time     = "2024-02-01T15:30:00Z"              // Updated value
  subject      = "Updated Test Scheduled Conference" // Updated value
  ews_item_id  = "updated-ews-item-id"               // Updated value
  ews_item_uid = "updated-ews-item-uid"              // Updated value

  depends_on = [
    pexip_infinity_conference.test-conference
  ]
}