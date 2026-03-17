/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_end_user" "tf-test-end-user" {
  primary_email_address = "tf-test-user@example.com"
  first_name            = "tf-test John"
  last_name             = "tf-test Doe"
  display_name          = "tf-test John Doe"
  telephone_number      = "+1234567890"
  mobile_number         = "+0987654321"
  title                 = "tf-test Software Engineer"
  department            = "tf-test Engineering"
  avatar_url            = "https://example.com/avatar.jpg"
  ms_exchange_guid      = "11111111-2222-3333-4444-555555555555"
  sync_tag              = "tf-test-sync-tag"
}
