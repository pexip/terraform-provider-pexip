/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider" "identity_provider-test" {
  name                                     = "identity_provider-test"
  description                              = "Test IdentityProvider"
  idp_type                                 = "saml"
  sso_url                                  = "https://example.com"
  idp_entity_id                            = "test-value"
  idp_public_key                           = "test-value"
  service_entity_id                        = "test-value"
  service_public_key                       = "test-value"
  service_private_key                      = "test-value"
  signature_algorithm                      = "rsa-sha256"
  digest_algorithm                         = "sha256"
  display_name_attribute_name              = "identity_provider-test"
  registration_alias_attribute_name        = "identity_provider-test"
  assertion_consumer_service_url           = "https://example.com"
  worker_fqdn_acs_urls                     = true
  disable_popup_flow                       = true
  oidc_flow                                = "authorization_code"
  oidc_client_id                           = "test-value"
  oidc_client_secret                       = "test-value"
  oidc_token_url                           = "https://example.com"
  oidc_user_info_url                       = "https://example.com"
  oidc_jwks_url                            = "https://example.com"
  oidc_token_endpoint_auth_scheme          = "client_secret_basic"
  oidc_token_signature_scheme              = "rs256"
  oidc_display_name_claim_name             = "identity_provider-test"
  oidc_registration_alias_claim_name       = "identity_provider-test"
  oidc_additional_scopes                   = "test-value"
  oidc_france_connect_required_eidas_level = "eidas1"
  attributes                               = "test-value"
}