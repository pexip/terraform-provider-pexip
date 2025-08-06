/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ldap_sync_field" "ldap_sync_field-test" {
  name                   = "ldap_sync_field-test"
  description            = "Updated Test LdapSyncField" // Updated description
  template_variable_name = "ldap_sync_field-test"
  is_binary              = false // Updated to false
}