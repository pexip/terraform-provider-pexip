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

resource "pexip_infinity_system_location" "test-location" {
  name        = "test-location"
  description = "Test Location"
  mtu         = 1460
}

resource "pexip_infinity_automatic_participant" "automatic-participant-test" {
  alias                 = "automatic-participant-updated"     // Updated value
  description           = "Updated Test AutomaticParticipant" // Updated description
  conference            = [pexip_infinity_conference.test-conference.id]
  protocol              = "h323"                      // Updated value
  call_type             = "audio"                     // Updated value
  role                  = "chair"                     // Updated value
  dtmf_sequence         = "456*"                      // Updated value
  keep_conference_alive = "keep_conference_alive" // Keep same value as API rejects end_conference_when_alone
  routing               = "manual"                    // Updated value
  system_location       = pexip_infinity_system_location.test-location.id
  streaming             = false // Updated to false
  remote_display_name   = "automatic_participant-test"
  presentation_url      = "https://updated.example.com" // Updated URL

  depends_on = [
    pexip_infinity_conference.test-conference,
    pexip_infinity_system_location.test-location
  ]
}