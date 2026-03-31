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

resource "pexip_infinity_mjx_integration" "test_min" {
  name             = "tf-test mjx-integration min"
  graph_deployment = pexip_infinity_mjx_graph_deployment.test_min.id
}

resource "pexip_infinity_mjx_integration" "test_full" {
  name             = "tf-test mjx-integration full"
  graph_deployment = pexip_infinity_mjx_graph_deployment.test_full.id
}

resource "pexip_infinity_mjx_meeting_processing_rule" "test" {
  name            = "tf-test mjx-meeting-processing-rule min"
  priority        = 1
  meeting_type    = "pexipinfinity"
  mjx_integration = pexip_infinity_mjx_integration.test_min.id
}
