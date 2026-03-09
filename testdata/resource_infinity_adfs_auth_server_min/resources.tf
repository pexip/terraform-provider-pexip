/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_adfs_auth_server" "tf-test-adfs-auth-server" {
  name                               = "tf-test-adfs-auth-server"
  client_id                          = "test-client-id-min"
  federation_service_name            = "adfs-min.example.com"
  federation_service_identifier      = "https://adfs-min.example.com/adfs/services/trust"
  relying_party_trust_identifier_url = "https://min.example.com"
}
