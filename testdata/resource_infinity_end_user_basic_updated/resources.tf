/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_end_user" "end-user-test" {
  primary_email_address = "updated@example.com"                    // Updated email
  first_name            = "Jane"                                   // Updated first name
  last_name             = "Smith"                                  // Updated last name
  display_name          = "Jane Smith"                             // Updated display name
  telephone_number      = "+1111111111"                            // Updated telephone
  mobile_number         = "+2222222222"                            // Updated mobile
  title                 = "Senior Engineer"                        // Updated title
  department            = "Product"                                // Updated department
  avatar_url            = "https://example.com/updated-avatar.jpg" // Updated avatar
  ms_exchange_guid      = "updated-guid-456"                       // Updated GUID
  sync_tag              = "updated-sync-tag"                       // Updated sync tag
}