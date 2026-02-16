/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_azure_tenant" "azure-tenant-test" {
  name        = "tf-test-azure-tenant-teams-proxy-full"
  description = "Test Azure Tenant for Teams Proxy"
  tenant_id   = "33333333-3333-3333-3333-333333333333"
}

resource "pexip_infinity_teams_proxy" "teams-proxy-test" {
  name                    = "tf-test-teams-proxy-full"
  description             = "Test Teams Proxy Full Configuration"
  address                 = "teams-proxy-full.pexvclab.com"
  port                    = 8443
  azure_tenant            = pexip_infinity_azure_tenant.azure-tenant-test.id
  min_number_of_instances = 3
  notifications_enabled   = true
  notifications_queue     = "Endpoint=sb://examplevmss.servicebus.windows.net/;SharedAccessKeyName=standard_access_policy;SharedAccessKey=testkey123="
  # eventhub_id is computed by the API from notifications_queue
}
