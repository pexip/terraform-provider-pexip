/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_teams_proxy" "teams-proxy-test" {
  name         = "test-teams-proxy"
  address      = "test-teams-proxy.dev.pexip.network"
  azure_tenant = "/api/admin/configuration/v1/azure_tenant/1/"
}