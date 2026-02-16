/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_azure_tenant" "azure-tenant-test" {
  name        = "tf-test-azure-tenant-teams-proxy-min"
  description = "Test Azure Tenant for Teams Proxy"
  tenant_id   = "12345678-1234-1234-1234-123456789012"
}

resource "pexip_infinity_teams_proxy" "teams-proxy-test" {
  name         = "tf-test-teams-proxy-min"
  address      = "teams-proxy-min.pexvclab.com"
  azure_tenant = pexip_infinity_azure_tenant.azure-tenant-test.id
}
