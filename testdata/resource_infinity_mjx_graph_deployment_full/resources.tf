/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_graph_deployment" "test" {
  name             = "tf-test mjx-graph-deployment full"
  description      = "Test MJX Graph deployment description"
  client_id        = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  client_secret    = "test-client-secret"
  oauth_token_url  = "https://login.microsoftonline.com/updated-tenant/oauth2/v2.0/token"
  graph_api_domain = "graph.microsoft.com"
  request_quota    = 500000
  disable_proxy    = true
}
