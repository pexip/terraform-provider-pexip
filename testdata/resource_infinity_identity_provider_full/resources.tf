/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  uuid = "988d1247-7997-46e9-a89a-5a148b5c5f29"
}

# Create TLS key and certificate for the IDP (simulating external identity provider)
resource "tls_private_key" "idp" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "idp" {
  private_key_pem = tls_private_key.idp.private_key_pem

  subject {
    common_name  = "tf-test-idp.example.com"
    organization = "Test IDP Organization"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

# Create TLS key and certificate for the service (Pexip)
resource "tls_private_key" "service" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "service" {
  private_key_pem = tls_private_key.service.private_key_pem

  subject {
    common_name  = "tf-test-service.pexip.example.com"
    organization = "Test Service Organization"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

# Create identity provider attributes for testing
resource "pexip_infinity_identity_provider_attribute" "attr1" {
  name        = "tf-test-displayName"
  description = "Test attribute for display name"
}

resource "pexip_infinity_identity_provider_attribute" "attr2" {
  name        = "tf-test-email"
  description = "Test attribute for email"
}

resource "pexip_infinity_identity_provider" "test" {
  # Required fields
  name                           = "tf-test Identity Provider full"
  uuid                           = local.uuid
  assertion_consumer_service_url = "https://test.example.com/oidcconsumer/${local.uuid}"

  # Optional basic fields
  description = "Full test Identity Provider with all fields"
  idp_type    = "oidc"

  # SAML specific fields
  sso_url                           = "https://idp.example.com/sso"
  idp_entity_id                     = "https://idp.example.com/entity"
  idp_public_key                    = tls_self_signed_cert.idp.cert_pem
  service_entity_id                 = "https://pexip.example.com/entity"
  service_public_key                = tls_self_signed_cert.service.cert_pem
  service_private_key               = tls_private_key.service.private_key_pem
  signature_algorithm               = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"
  digest_algorithm                  = "http://www.w3.org/2001/04/xmldsig-more#sha384"
  display_name_attribute_name       = "displayName"
  registration_alias_attribute_name = "userPrincipalName"

  # Additional assertion consumer service URLs
  assertion_consumer_service_url2  = "https://test2.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url3  = "https://test3.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url4  = "https://test4.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url5  = "https://test5.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url6  = "https://test6.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url7  = "https://test7.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url8  = "https://test8.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url9  = "https://test9.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"
  assertion_consumer_service_url10 = "https://test10.example.com/oidcconsumer/988d1247-7997-46e9-a89a-5a148b5c5f29"

  # Worker and popup settings
  worker_fqdn_acs_urls = true
  disable_popup_flow   = true

  # OIDC specific fields
  oidc_flow                                = "implicit"
  oidc_client_id                           = "test-client-id-12345"
  oidc_client_secret                       = "test-client-secret-67890"
  oidc_token_url                           = "https://idp.example.com/oauth2/token"
  oidc_user_info_url                       = "https://idp.example.com/oauth2/userinfo"
  oidc_jwks_url                            = "https://idp.example.com/.well-known/jwks.json"
  oidc_token_endpoint_auth_scheme          = "client_secret_basic"
  oidc_token_signature_scheme              = "hs256"
  oidc_display_name_claim_name             = "full_name"
  oidc_registration_alias_claim_name       = "preferred_username"
  oidc_additional_scopes                   = "profile email phone address"
  oidc_france_connect_required_eidas_level = "eidas3"

  # Attributes
  attributes = [
    pexip_infinity_identity_provider_attribute.attr1.id,
    pexip_infinity_identity_provider_attribute.attr2.id,
  ]
}
