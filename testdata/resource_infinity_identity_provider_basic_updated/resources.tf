/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider" "identity_provider-test" {
  name                                     = "identity_provider-test"
  description                              = "Updated Test IdentityProvider" // Updated description
  idp_type                                 = "oidc"                          // Updated value
  sso_url                                  = "https://updated.example.com"   // Updated URL
  idp_entity_id                            = "updated-value"                 // Updated value
  idp_public_key                           = "updated-value"                 // Updated value
  service_entity_id                        = "updated-value"                 // Updated value
  service_public_key                       = "updated-value"                 // Updated value
  service_private_key                      = "updated-value"                 // Updated value
  signature_algorithm                      = "rsa-sha1"                      // Updated value
  digest_algorithm                         = "sha1"                          // Updated value
  display_name_attribute_name              = "identity_provider-test"
  registration_alias_attribute_name        = "identity_provider-test"
  assertion_consumer_service_url           = "https://updated.example.com" // Updated URL
  worker_fqdn_acs_urls                     = false                         // Updated to false
  disable_popup_flow                       = false                         // Updated to false
  oidc_flow                                = "implicit"                    // Updated value
  oidc_client_id                           = "updated-value"               // Updated value
  oidc_client_secret                       = "updated-value"               // Updated value
  oidc_token_url                           = "https://updated.example.com" // Updated URL
  oidc_user_info_url                       = "https://updated.example.com" // Updated URL
  oidc_jwks_url                            = "https://updated.example.com" // Updated URL
  oidc_token_endpoint_auth_scheme          = "client_secret_post"          // Updated value
  oidc_token_signature_scheme              = "hs256"                       // Updated value
  oidc_display_name_claim_name             = "identity_provider-test"
  oidc_registration_alias_claim_name       = "identity_provider-test"
  oidc_additional_scopes                   = "updated-value" // Updated value
  oidc_france_connect_required_eidas_level = "eidas2"        // Updated value
  attributes                               = "updated-value" // Updated value
}