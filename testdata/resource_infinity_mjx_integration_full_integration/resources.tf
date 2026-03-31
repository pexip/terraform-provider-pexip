/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_graph_deployment" "test_min" {
  name            = "tf-test mjx-graph-deployment min"
  client_id       = "12345678-1234-1234-1234-123456789012"
  oauth_token_url = "https://login.microsoftonline.com/test-tenant/oauth2/v2.0/token"
}

resource "pexip_infinity_mjx_graph_deployment" "test_full" {
  name            = "tf-test mjx-graph-deployment full"
  client_id       = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  oauth_token_url = "https://login.microsoftonline.com/updated-tenant/oauth2/v2.0/token"
}

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
  graph_deployment               = pexip_infinity_mjx_graph_deployment.test_full.id
  process_alias_private_meetings = false
  replace_empty_subject          = false
  replace_subject_type           = "ALL"
  replace_subject_template       = "Meeting: {{ subject }}"
  use_webex                      = false
  webex_api_domain               = "custom.webexapis.com"
}
