/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "INFINITY_GMS_GW_TOKEN_CERT2" {
  description = "The certificate chain for the GMS gateway token"
  type        = string
}

variable "INFINITY_GMS_GW_TOKEN_KEY2" {
  description = "The private key for the GMS gateway token"
  type        = string
}

resource "pexip_infinity_gms_gateway_token" "test" {
  certificate = var.INFINITY_GMS_GW_TOKEN_CERT2
  private_key  = var.INFINITY_GMS_GW_TOKEN_KEY2
}
