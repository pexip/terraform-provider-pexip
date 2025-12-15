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
  start_time   = "2024-01-01T10:00:00Z"
  end_time     = "2024-01-01T11:00:00Z"
  subject      = "Test Scheduled Conference"
  ews_item_id  = "test-ews-item-id"
  ews_item_uid = "test-ews-item-uid"

  depends_on = [
    pexip_infinity_conference.test-conference
  ]
}