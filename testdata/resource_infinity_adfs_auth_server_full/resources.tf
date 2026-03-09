/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_adfs_auth_server" "tf-test-adfs-auth-server" {
  name                               = "tf-test-adfs-auth-server"
  description                        = "Full test configuration for ADFS Auth Server"
  client_id                          = "test-client-id-full"
  federation_service_name            = "adfs-full.example.com"
  federation_service_identifier      = "https://adfs-full.example.com/adfs/services/trust"
  relying_party_trust_identifier_url = "https://full.example.com"
}
