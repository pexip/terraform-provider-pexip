/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider" "identity_provider-test" {
  name                               = "identity_provider-test"
  description                        = "Test IdentityProvider"
  idp_type                           = "saml"
  sso_url                            = "https://example.com/sso"
  idp_entity_id                      = "https://example.com/entity"
  service_entity_id                  = "https://pexip.example.com/entity"
  signature_algorithm                = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
  digest_algorithm                   = "http://www.w3.org/2001/04/xmlenc#sha256"
  display_name_attribute_name        = "displayName"
  registration_alias_attribute_name  = "email"
  worker_fqdn_acs_urls               = true
  disable_popup_flow                 = true
  oidc_flow                          = "code"
  oidc_token_endpoint_auth_scheme    = "client_secret_basic"
  oidc_token_signature_scheme        = "rs256"
  oidc_france_connect_required_eidas_level = "eidas1"
}