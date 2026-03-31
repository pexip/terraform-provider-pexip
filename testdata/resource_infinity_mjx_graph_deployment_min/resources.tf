/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_mjx_graph_deployment" "test" {
  name            = "tf-test mjx-graph-deployment min"
  client_id       = "12345678-1234-1234-1234-123456789012"
  oauth_token_url = "https://login.microsoftonline.com/test-tenant/oauth2/v2.0/token"
}
