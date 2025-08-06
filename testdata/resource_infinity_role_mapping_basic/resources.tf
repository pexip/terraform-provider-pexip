/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_role_mapping" "role_mapping-test" {
  name   = "role_mapping-test"
  source = "saml_attribute"
  value  = "test-value"
}