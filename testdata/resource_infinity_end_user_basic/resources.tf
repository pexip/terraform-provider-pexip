/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_end_user" "end-user-test" {
  primary_email_address = "user@example.com"
  first_name            = "John"
  last_name             = "Doe"
  display_name          = "John Doe"
  telephone_number      = "+1234567890"
  mobile_number         = "+0987654321"
  title                 = "Software Engineer"
  department            = "Engineering"
  avatar_url            = "https://example.com/avatar.jpg"
  ms_exchange_guid      = "ms-guid-123"
  sync_tag              = "sync-tag"
}