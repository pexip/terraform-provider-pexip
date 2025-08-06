/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_azure_tenant" "azure_tenant-test" {
  name        = "azure_tenant-test"
  description = "Updated Test AzureTenant" // Updated description
  tenant_id   = "updated-value"            // Updated value
}