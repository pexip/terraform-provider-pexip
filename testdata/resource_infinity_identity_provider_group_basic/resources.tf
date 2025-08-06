/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_identity_provider_group" "identity_provider_group-test" {
  name        = "identity_provider_group-test"
  description = "Test IdentityProviderGroup"
}