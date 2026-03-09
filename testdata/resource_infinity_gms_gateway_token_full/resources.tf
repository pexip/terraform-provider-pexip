/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

variable "infinity_gms_gw_token_cert2" {
  description = "The certificate chain for the GMS gateway token"
  type        = string
}

variable "infinity_gms_gw_token_key2" {
  description = "The private key for the GMS gateway token"
  type        = string
}

resource "pexip_infinity_gms_gateway_token" "test" {
  certificate = var.infinity_gms_gw_token_cert2
  private_key = var.infinity_gms_gw_token_key2
}
