/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "random_uuid4" "example" {
}

locals {
  uuid = "988d1247-7997-46e9-a89a-5a148b5c5f29"
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
  uuid = local.uuid
  assertion_consumer_service_url = "https://test.com/samlconsumer/${local.uuid}"
  # Note: attributes not referenced here, so the identity_provider will clear them
}