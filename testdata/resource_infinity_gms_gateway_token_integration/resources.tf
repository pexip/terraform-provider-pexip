/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "gms_certificate" {
  description = "GMS gateway token certificate chain"
  type        = string
  sensitive   = true
}

variable "gms_private_key" {
  description = "GMS gateway token private key"
  type        = string
  sensitive   = true
}

resource "pexip_infinity_gms_gateway_token" "gms_gateway_token_integration_test" {
  certificate = var.gms_certificate
  private_key = var.gms_private_key
}
