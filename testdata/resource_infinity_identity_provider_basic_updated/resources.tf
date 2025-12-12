/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider" "identity_provider-test" {
  name                                     = "identity_provider-test-updated"
  description                              = "Updated Test IdentityProvider"
  idp_type                                 = "oidc"
  sso_url                                  = "https://updated.example.com/sso"
  idp_entity_id                            = "https://updated.example.com/entity"
  service_entity_id                        = "https://pexip-updated.example.com/entity"
  signature_algorithm                      = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
  digest_algorithm                         = "http://www.w3.org/2001/04/xmlenc#sha1"
  display_name_attribute_name              = "name"
  registration_alias_attribute_name        = "username"
  worker_fqdn_acs_urls                     = false
  disable_popup_flow                       = false
  oidc_flow                                = "implicit"
  oidc_client_id                           = "updated-client-id"
  oidc_client_secret                       = "updated-client-secret"
  oidc_token_url                           = "https://updated.example.com/token"
  oidc_user_info_url                       = "https://updated.example.com/userinfo"
  oidc_jwks_url                            = "https://updated.example.com/jwks"
  oidc_token_endpoint_auth_scheme          = "client_secret_post"
  oidc_token_signature_scheme              = "hs256"
  oidc_display_name_claim_name             = "name"
  oidc_registration_alias_claim_name       = "preferred_username"
  oidc_additional_scopes                   = "profile email"
  oidc_france_connect_required_eidas_level = "eidas2"
}