/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "tls_private_key" "test" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test" {
  private_key_pem = tls_private_key.test.private_key_pem

  subject {
    common_name  = "pexip.network.updated"
    organization = "Pexip Org"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "pexip_infinity_tls_certificate" "tls-cert-test" {
  certificate = tls_self_signed_cert.test.cert_pem
  private_key = tls_private_key.test.private_key_pem
}
