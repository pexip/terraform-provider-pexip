/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "random_uuid4" "example" {
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
  name = "tf-test Identity Provider min"
  uuid = random_uuid4.example.result
  assertion_consumer_service_url = "https://test.com/samlconsumer/${random_uuid4.example.result}"
  # Note: attributes not referenced here, so the identity_provider will clear them
}