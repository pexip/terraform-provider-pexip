/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_gateway_routing_rule" "tf-test-gateway-routing-rule" {
  name                            = "tf-test-gateway-routing-rule"
  description                     = "tf-test Gateway Routing Rule Description"
  match_string                    = ".*@example.com"
  priority                        = 100
  enable                          = false
  match_string_full               = true
  replace_string                  = "replaced@example.com"
  tag                             = "tf-test-tag"

  # Call settings
  called_device_type              = "registration"
  outgoing_protocol               = "h323"
  call_type                       = "audio"
  crypto_mode                     = "on"

  # Audio/Video settings
  denoise_audio                   = false
  max_pixels_per_second           = "fullhd"
  max_callrate_in                 = 2048
  max_callrate_out                = 4096

  # Matching rules
  match_incoming_calls            = true
  match_outgoing_calls            = true
  match_incoming_sip              = false
  match_incoming_h323             = false
  match_incoming_mssip            = false
  match_incoming_webrtc           = false
  match_incoming_teams            = true
  match_incoming_only_if_registered = false

  # Features
  enable_participant_avatar_lookup = "yes"
  live_captions_enabled            = "yes"
  treat_as_trusted                 = true
}
