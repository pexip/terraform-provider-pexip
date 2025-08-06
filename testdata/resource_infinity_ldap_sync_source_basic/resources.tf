/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

resource "pexip_infinity_ldap_sync_source" "ldap_sync_source-test" {
  name                    = "ldap_sync_source-test"
  description             = "Test LdapSyncSource"
  ldap_server             = "test-value"
  ldap_base_dn            = "test-value"
  ldap_bind_username      = "ldap_sync_source-test"
  ldap_bind_password      = "test-value"
  ldap_use_global_catalog = true
  ldap_permit_no_tls      = true
}