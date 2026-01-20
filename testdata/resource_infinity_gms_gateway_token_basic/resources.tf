/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_gms_gateway_token" "gms_gateway_token_test" {
  certificate = <<-EOT
-----BEGIN CERTIFICATE-----
server-cert
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
intermediate-cert
-----END CERTIFICATE-----
EOT
  private_key = "test-private-key-data"
}
