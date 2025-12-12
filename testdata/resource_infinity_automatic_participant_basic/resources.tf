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
  alias                 = "automatic-participant-test"
  description           = "Test AutomaticParticipant"
  conference            = [pexip_infinity_conference.test-conference.id]
  protocol              = "sip"
  call_type             = "video"
  role                  = "guest"
  dtmf_sequence         = "123#"
  keep_conference_alive = "keep_conference_alive"
  routing               = "manual"
  system_location       = pexip_infinity_system_location.test-location.id
  streaming             = true
  remote_display_name   = "automatic_participant-test"
  presentation_url      = "https://example.com"

  depends_on = [
    pexip_infinity_conference.test-conference,
    pexip_infinity_system_location.test-location
  ]
}