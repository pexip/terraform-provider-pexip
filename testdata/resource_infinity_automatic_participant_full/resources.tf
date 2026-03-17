/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_conference" "test-adp-conference" {
  name         = "tf-test-adp-conference"
  service_type = "conference"
}

resource "pexip_infinity_system_location" "test-adp-location" {
  name = "tf-test-adp-location"
}

resource "pexip_infinity_automatic_participant" "automatic-participant-test" {
  alias                 = "tf-test-adp"
  description           = "Test AutomaticParticipant Full Config"
  conference            = [pexip_infinity_conference.test-adp-conference.id]
  protocol              = "h323"
  call_type             = "audio"
  role                  = "chair"
  dtmf_sequence         = "456*"
  keep_conference_alive = "keep_conference_alive"
  routing               = "routing_rule"
  system_location       = pexip_infinity_system_location.test-adp-location.id
  streaming             = true
  remote_display_name   = "tf-test-adp-display"
  presentation_url      = "https://example.com/presentation"
}