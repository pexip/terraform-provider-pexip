/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_azure_tenant" "azure-tenant-test" {
  name        = "tf-test azure-tenant-teams-proxy scheduled scaling"
  description = "Test Azure Tenant for Scheduled Scaling"
  tenant_id   = "44444444-4444-4444-4444-444444444445"
}

resource "pexip_infinity_teams_proxy" "teams-proxy-test" {
  name         = "tf-test teams-proxy scheduled scaling"
  address      = "teams-proxy-min.pexvclab.com"
  azure_tenant = pexip_infinity_azure_tenant.azure-tenant-test.id
}

resource "pexip_infinity_scheduled_scaling" "test" {
  policy_name         = "tf-test scheduled scaling min"
  policy_type         = "TeamsConnectorScaling"
  resource_identifier = pexip_infinity_teams_proxy.teams-proxy-test.name
  start_date          = "2024-01-01"
  time_from           = "09:00:00"
  time_to             = "17:00:00"
  mon                 = true
}