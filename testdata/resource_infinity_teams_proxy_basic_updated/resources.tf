/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_teams_proxy" "teams-proxy-test" {
  name                    = "test-teams-proxy-updated"
  description             = "Test Teams Proxy Updated"
  address                 = "updated-test-teams-proxy.dev.pexip.network"
  port                    = 8443
  azure_tenant            = "/api/admin/configuration/v1/azure_tenant/1/"
  eventhub_id             = "updated-test-eventhub-id"
  min_number_of_instances = 0
  notifications_queue     = "updated-test-notifications-queue"
  notifications_enabled   = true
}