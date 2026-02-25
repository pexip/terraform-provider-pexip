/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

# Create TLS private key for OIDC authentication
# Kept to avoid deletion when switching between full and min configs
resource "tls_private_key" "oidc" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "pexip_infinity_authentication" "authentication-test" {
  # Note: oidc_private_key not referenced here, so the authentication will clear it
}