/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_adfs_auth_server" "adfs_auth_server-test" {
  name                               = "adfs_auth_server-test"
  description                        = "Test ADFSAuthServer"
  client_id                          = "test-client-id-12345"
  federation_service_name            = "adfs.example.com"
  federation_service_identifier      = "https://adfs.example.com/adfs/services/trust"
  relying_party_trust_identifier_url = "https://example.com"
}