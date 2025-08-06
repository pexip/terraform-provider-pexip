/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_scheduled_conference" "scheduled_conference-test" {
  conference           = "test-value"
  start_time           = "2024-01-01T10:00:00Z"
  end_time             = "2024-01-01T11:00:00Z"
  subject              = "test-value"
  ews_item_id          = "test-value"
  ews_item_uid         = "test-value"
  recurring_conference = "test-value"
  scheduled_alias      = "test-value"
}