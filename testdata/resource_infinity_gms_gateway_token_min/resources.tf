/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "infinity_gms_gw_token_cert" {
  description = "The certificate for the GMS gateway token"
  type        = string
}

variable "infinity_gms_gw_token_key" {
  description = "The private key for the GMS gateway token"
  type        = string
}

resource "pexip_infinity_gms_gateway_token" "test" {
  certificate = var.infinity_gms_gw_token_cert
  private_key  = var.infinity_gms_gw_token_key
}
