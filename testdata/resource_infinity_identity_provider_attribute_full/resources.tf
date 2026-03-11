/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider_attribute" "tf-test-identity-provider-attribute" {
  name        = "tf-test-identity-provider-attribute"
  description = "Test Identity Provider Attribute Description"
}
