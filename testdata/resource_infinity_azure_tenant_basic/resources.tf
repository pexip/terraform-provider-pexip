/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_azure_tenant" "azure_tenant-test" {
  name        = "azure_tenant-test"
  description = "Test AzureTenant"
  tenant_id   = "test-value"
}