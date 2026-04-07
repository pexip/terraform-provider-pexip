/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_global_configuration" "config" {
  enable_breakout_rooms = true
}

resource "pexip_infinity_automatic_participant" "tf-test-participant1" {
  alias = "tf-test-participant1@example.com"
}

resource "pexip_infinity_automatic_participant" "tf-test-participant2" {
  alias = "tf-test-participant2@example.com"
}

resource "pexip_infinity_conference" "tf-test-conference" {
  name                                = "tf-test-conference"
  description                         = "Full test configuration for conference"
  service_type                        = "conference"
  pin                                 = "123456"
  guest_pin                           = "654321"
  tag                                 = "tf-test-tag"
  allow_guests                        = true
  automatic_participants              = [pexip_infinity_automatic_participant.tf-test-participant1.id, pexip_infinity_automatic_participant.tf-test-participant2.id]
  breakout_rooms                      = true
  call_type                           = "video-only"
  crypto_mode                         = "besteffort"
  denoise_enabled                     = true
  direct_media                        = "best_effort"
  direct_media_notification_duration  = 5
  enable_active_speaker_indication    = true
  enable_chat                         = "yes"
  enable_overlay_text                 = true
  force_presenter_into_main           = true
  guests_can_present                  = true
  guests_can_see_guests               = "always"
  host_view                           = "two_mains_twentyone_pips"
  live_captions_enabled               = "yes"
  match_string                        = "^[0-9]+$"
  max_callrate_in                     = 4096
  max_callrate_out                    = 2048
  max_pixels_per_second               = "fullhd"
  mute_all_guests                     = true
  non_idp_participants                = "allow_if_trusted"
  on_completion                       = "{\"disconnect\": true}"
  participant_limit                   = 50
  post_match_string                   = "^test"
  post_replace_string                 = "new-test"
  primary_owner_email_address         = "owner@example.com"
  replace_string                      = "replaced"
  softmute_enabled                    = true
  sync_tag                            = "sync-123"
  two_stage_dial_type                 = "regular"

  depends_on = [
    pexip_infinity_global_configuration.config
  ]
}

resource "pexip_infinity_conference_alias" "tf-test-alias1" {
  alias       = "tf-test-alias1"
  description = "Test alias 1"
  conference  = pexip_infinity_conference.tf-test-conference.id
}

resource "pexip_infinity_conference_alias" "tf-test-alias2" {
  alias       = "tf-test-alias2"
  description = "Test alias 2"
  conference  = pexip_infinity_conference.tf-test-conference.id
}
