/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_azure_tenant" "azure_tenant-test" {
  name        = "tf-test-azure-tenant-min"
  description = "" // Explicitly clear description
  tenant_id   = "87654321-4321-4321-4321-210987654321" // Updated UUID
}