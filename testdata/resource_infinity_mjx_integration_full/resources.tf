/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_integration" "test" {
  name                           = "tf-test mjx-integration full"
  description                    = "Test MJX integration description"
  display_upcoming_meetings      = 14
  enable_non_video_meetings      = false
  enable_private_meetings        = false
  end_buffer                     = 10
  start_buffer                   = 10
  ep_username                    = "ep-user@example.com"
  ep_password                    = "ep-password-test"
  ep_use_https                   = false
  ep_verify_certificate          = true
  graph_deployment               = "/api/admin/configuration/v1/mjx_graph_deployment/2/"
  process_alias_private_meetings = false
  replace_empty_subject          = false
  replace_subject_type           = "ALL"
  replace_subject_template       = "Meeting: {{ subject }}"
  use_webex                      = true
  webex_api_domain               = "custom.webexapis.com"
  webex_client_id                = "webex-client-id-full"
  webex_client_secret            = "webex-secret-full"
  webex_redirect_uri             = "https://pexip.example.com/admin/platform/mjxintegration/oauth_redirect/"
}
