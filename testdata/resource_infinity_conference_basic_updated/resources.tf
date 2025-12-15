/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "conference-test" {
  name         = "conference-test"
  description  = "Updated Test Conference" // Updated description
  service_type = "conference"              // Keep same service type (not updatable)
  pin          = "9876"                    // Updated PIN
  guest_pin    = "4321"                    // Updated guest PIN
  allow_guests = true                      // Updated to true
  tag          = "updated-tag"             // Updated tag
}