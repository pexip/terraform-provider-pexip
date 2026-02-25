/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "infinity_licence_key2" {
  description = "The entitlement ID for the licence activation."
  default = "test-value"
  type        = string
  
}

resource "pexip_infinity_licence" "licence-test" {
  entitlement_id = var.infinity_licence_key2
}