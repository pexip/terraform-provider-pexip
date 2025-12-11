/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_automatic_participant" "automatic-participant-test" {
  alias                 = "automatic-participant-test"
  description           = "Test AutomaticParticipant"
  conference            = "/api/admin/configuration/v1/conference/1/"
  protocol              = "sip"
  call_type             = "video"
  role                  = "guest"
  dtmf_sequence         = "123#"
  keep_conference_alive = "keep_conference_alive"
  routing               = "auto"
  system_location       = "/api/admin/configuration/v1/system_location/1/"
  streaming             = true
  remote_display_name   = "automatic_participant-test"
  presentation_url      = "https://example.com"
}