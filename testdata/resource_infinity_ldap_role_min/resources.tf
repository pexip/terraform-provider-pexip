/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ldap_role" "ldap_role-test" {
  name          = "tf-test min"
  ldap_group_dn = "testmin"
}