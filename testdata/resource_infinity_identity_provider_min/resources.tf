/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

locals {
  uuid = "988d1247-7997-46e9-a89a-5a148b5c5f29"
}

# Create TLS key and certificate for the IDP (simulating external identity provider)
# Kept to avoid deletion when switching between full and min configs
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
# Kept to avoid deletion when switching between full and min configs
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

# Keep attribute resources to avoid deletion before identity_provider is updated
resource "pexip_infinity_identity_provider_attribute" "attr1" {
  name        = "tf-test-displayName"
  description = "Test attribute for display name"
}

resource "pexip_infinity_identity_provider_attribute" "attr2" {
  name        = "tf-test-email"
  description = "Test attribute for email"
}

resource "pexip_infinity_identity_provider" "test" {
  name                           = "tf-test Identity Provider min"
  uuid                           = local.uuid
  assertion_consumer_service_url = "https://test.com/samlconsumer/${local.uuid}"
  # Note: attributes and TLS keys not referenced here, so the identity_provider will clear them
}