/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_global_configuration" "global_configuration-test" {
  enable_webrtc                 = true
  enable_sip                    = true
  enable_h323                   = true
  enable_rtmp                   = true
  crypto_mode                   = "besteffort"
  max_pixels_per_second         = "hd"
  bursting_enabled              = true
  cloud_provider                = "AWS"
  aws_access_key                = "test-key"
  aws_secret_key                = "test-secret"
  azure_client_id               = "test-client"
  azure_secret                  = "test-secret"
  enable_analytics              = true
  media_ports_start             = 40000
  media_ports_end               = 40100
  signalling_ports_start        = 5060
  signalling_ports_end          = 5070
  guests_only_timeout           = 300
  waiting_for_chair_timeout     = 600
}