/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_adfs_auth_server" "adfs_auth_server-test" {
  name                               = "adfs_auth_server-test"
  description                        = "Updated Test ADFSAuthServer"                          // Updated description
  client_id                          = "updated-client-id-67890"                              // Updated value
  federation_service_name            = "adfs-updated.example.com"                             // Updated FQDN
  federation_service_identifier      = "https://adfs-updated.example.com/adfs/services/trust" // Updated URL
  relying_party_trust_identifier_url = "https://updated.example.com"                          // Updated URL
}